# System Design Questions — Linux Update Systems (Tesla)

These are design/discussion questions modeled after what you might be asked in a
Tesla OTA update systems internship interview. For each question, think through
your answer for 5-10 minutes before reading the hints.

Practice articulating your answers out loud — interviewers care about your
thought process as much as the answer itself.

---

## 1. Design an OTA Update Pipeline

> You're tasked with designing a system that delivers software updates to 2 million
> vehicles in the field. Each vehicle has an infotainment system (Linux-based) and
> several microcontrollers (body controller, battery management, etc.).
>
> Walk through how you'd design the end-to-end pipeline from "engineer merges code"
> to "vehicle is running the new software."

**Things to cover:**
- Build artifacts and versioning
- How updates are packaged (full image vs delta)
- Server-side: how do you stage, roll out, and track updates?
- Client-side: how does the vehicle download, verify, and apply an update?
- How do you handle the fact that different components update differently?
- What happens if the vehicle loses connectivity mid-download?

<details>
<summary>Hints / discussion points</summary>

- **Staged rollout**: canary (internal fleet) → early access → % ramp → full fleet.
  Think about how to pause/halt a rollout if error rates spike.
- **Delta vs full**: Delta saves bandwidth but adds complexity. Consider bsdiff/xdelta.
  Full images are simpler and easier to verify. Tesla likely uses both depending on component.
- **Manifest file**: A signed JSON/protobuf describing which components need updating,
  expected versions, checksums, and ordering constraints.
- **Download resumption**: Range requests over HTTPS. Store partial downloads and
  verify chunks incrementally.
- **Component orchestration**: Some components must update in order (e.g., gateway
  firmware before body controllers). Model this as a DAG.
</details>

---

## 2. Atomic Updates and Rollback

> A vehicle starts an update to its infotainment Linux system. Power could be cut
> at any point (user unplugs, 12V battery dies). Design the update mechanism so
> the system is **never left in an unbootable state**.

**Things to cover:**
- What does "atomic" mean in this context?
- A/B partition scheme vs in-place updates
- How does the bootloader know which partition to boot?
- What triggers a rollback?
- How do you handle the filesystem, kernel, and initramfs as a unit?

<details>
<summary>Hints / discussion points</summary>

- **A/B partitioning**: Two root filesystem slots. Write the update to the inactive
  slot, flip a flag, reboot. If the new slot fails health checks, revert the flag.
  This is the industry standard (Android uses this, ChromeOS, etc.).
- **Bootloader integration**: U-Boot or GRUB reads a "boot_slot" variable from
  a small metadata partition. After N failed boots, it flips back.
- **dm-verity**: Read-only verified root filesystem. If any block is corrupt,
  the read fails and triggers rollback.
- **Health check**: After booting the new slot, a watchdog service must confirm
  success within X seconds or the system reboots into the old slot.
- **Separate concerns**: Kernel + initramfs + rootfs are bundled per slot.
  User data lives on a separate partition that survives updates.
</details>

---

## 3. Secure Update Delivery

> An attacker intercepts the network connection between Tesla's servers and a vehicle.
> How do you ensure the vehicle **only installs authentic, unmodified software**
> from Tesla?

**Things to cover:**
- Code signing and verification
- Chain of trust from Tesla's build system to the vehicle
- What keys live where? How are they protected?
- What's the threat model? (MITM, compromised server, stolen vehicle)
- How do you handle key rotation or revocation?

<details>
<summary>Hints / discussion points</summary>

- **Asymmetric signing**: Tesla signs update packages with a private key (stored in
  HSM). Vehicles verify with the corresponding public key baked into firmware/bootloader.
- **Chain of trust**: Secure Boot → bootloader verifies kernel → kernel verifies
  rootfs (dm-verity) → update agent verifies downloaded packages.
- **TLS is not enough**: TLS protects the transport, but the artifact itself must be
  signed. A compromised CDN or MITM proxy could serve bad files over valid TLS.
- **Key rotation**: Embed multiple public keys or a certificate chain. New keys can
  be delivered via a signed update before the old key is retired.
- **Metadata signing**: The update manifest is signed separately from the payload.
  Vehicle checks manifest signature before downloading anything.
- **Replay protection**: Include a monotonic version counter. Vehicle refuses to
  install anything older than current version.
</details>

---

## 4. Updating Microcontrollers via a Linux Host

> The vehicle's Linux infotainment system needs to update firmware on a body
> controller connected via CAN bus, and a camera module connected via SPI.
> Design the update flow.

**Things to cover:**
- How does the Linux host communicate the update to peripheral MCUs?
- What protocols would you use (CAN/UDS, SPI, I2C)?
- How do you make MCU updates atomic/safe?
- What if the MCU has only 256KB of flash — can you do A/B?
- Error handling: what if the MCU update fails halfway?

<details>
<summary>Hints / discussion points</summary>

- **UDS over CAN**: Unified Diagnostic Services is standard in automotive for
  flashing ECUs. Sequence: request download → transfer data → request transfer exit
  → reset ECU.
- **Bootloader on MCU**: The MCU has a small bootloader that can receive firmware
  over CAN/SPI. If the main application is corrupt, the bootloader can re-flash.
- **256KB flash constraint**: Can't do full A/B. Instead, use a bootloader +
  application partition. Bootloader verifies application CRC on boot. If invalid,
  it enters "recovery mode" waiting for a reflash.
- **Chunked transfer with CRC**: Send firmware in chunks, each with a CRC.
  MCU acknowledges each chunk. Final step: CRC of entire image, then commit.
- **Ordering matters**: If the gateway ECU needs to relay CAN messages to body
  controllers, update the gateway last (or ensure backward compatibility).
- **Timeout and retry**: If a CAN transfer stalls, the Linux host retries.
  After N failures, mark the update as failed and report to the server.
</details>

---

## 5. Fleet Rollout Strategy and Metrics

> You've built the update system. Now product wants to roll out a new firmware
> version to all vehicles. How do you design the rollout process and what
> metrics would you monitor?

**Things to cover:**
- Rollout phases (internal → canary → staged → full)
- What metrics indicate a bad update? How do you detect them?
- Automatic rollback criteria
- How do you handle vehicles that are offline for days/weeks?
- Dashboard design: what does the ops team need to see?

<details>
<summary>Hints / discussion points</summary>

- **Phased rollout**:
  - Internal/dogfood fleet (employees): catch obvious issues.
  - Canary (1%): instrument heavily, compare error rates to baseline.
  - Staged (10% → 25% → 50% → 100%): watch for long-tail issues.
- **Key metrics**:
  - Update success/failure rate per component
  - Boot success rate post-update (are vehicles getting stuck?)
  - Rollback rate
  - Download retry rate and average download time
  - Time from "update available" to "update complete"
  - Error codes distribution
- **Automatic halt**: If failure rate exceeds X% in any cohort, pause rollout.
  Alert on-call engineer. This is non-negotiable for safety-critical systems.
- **Offline vehicles**: Updates remain "pending." Vehicle checks in when it
  connects to WiFi. May need to skip versions (1.0 → 1.3 directly).
  Ensure the update system supports non-sequential upgrades.
- **A/B testing**: Sometimes you want to compare two firmware versions.
  Assign vehicles to groups and compare metrics.
</details>

---

## 6. Delta Updates for Bandwidth Efficiency

> Vehicles primarily update over WiFi, but some use cellular (LTE). A full
> rootfs image is 4GB. Design a delta update system that minimizes download size.

**Things to cover:**
- Binary diff algorithms (bsdiff, xdelta, etc.)
- How do you generate deltas on the server side?
- What version pairs do you need deltas for?
- Client-side: applying the delta, verification
- Fallback: when should you fall back to a full image?

<details>
<summary>Hints / discussion points</summary>

- **Binary diff**: bsdiff produces very compact patches for binaries. xdelta3 is
  faster but sometimes larger. Google's Puffin works well on compressed data.
- **Delta matrix problem**: If you have versions 1.0, 1.1, 1.2, 1.3 in the field,
  do you generate deltas for every pair? That's O(n²). In practice, generate
  deltas from the N most common versions to the target version.
- **Server-side**: Build pipeline generates full image + deltas from top-N source
  versions. Store on CDN. Manifest tells the vehicle which delta to download
  based on its current version.
- **Client-side apply**: Download delta, apply to current partition, write to
  inactive partition. Verify checksum of result matches expected hash.
- **Fallback**: If no delta exists for the vehicle's current version, or if delta
  application fails, download the full image. Must always work.
- **Compression**: Delta is already compact; additional gzip/zstd compression
  on top helps. Stream decompression on the vehicle to save disk space.
</details>

---

## 7. Concurrent Update Orchestration in Go

> The vehicle needs to update 8 components: infotainment OS, autopilot firmware,
> 4 body controllers, a gateway ECU, and navigation maps. Some have dependencies
> (gateway must update before body controllers). Design the orchestration
> layer in Go.

**Things to cover:**
- Data structure for representing update dependencies (DAG)
- How do you execute independent updates in parallel?
- How do you handle partial failures? (3 of 8 succeed, 2 fail, 3 not started)
- How does this interact with the A/B or rollback mechanism?
- Go-specific: what concurrency primitives would you use?

<details>
<summary>Hints / discussion points</summary>

- **DAG representation**: Each component is a node. Edges represent "must complete
  before." Use adjacency list. Topological sort gives execution order.
- **Go implementation**:
  - Each component update is a goroutine.
  - Use channels or `sync.WaitGroup` to signal completion.
  - A component goroutine blocks until all its dependencies have signaled success.
  - Use `context.Context` for cancellation on failure.
- **Parallel execution**: Components with no unmet dependencies run concurrently.
  E.g., nav maps and infotainment OS can update in parallel if independent.
- **Failure modes**:
  - If a critical component fails, cancel all dependents.
  - If a non-critical component fails, continue with others, report partial success.
  - Persist state: "component X succeeded, Y failed" so a retry doesn't redo X.
- **Rollback coordination**: If the infotainment update succeeds but a body
  controller update fails, do you roll back the infotainment? Depends on
  compatibility. The manifest should declare compatibility constraints.
- **State machine per component**: Pending → Downloading → Applying → Verifying →
  Complete | Failed. Persist state to disk in case of power loss.
</details>

---

## 8. Debugging a Fleet-Wide Update Failure

> After rolling out a new update, 5% of vehicles report "update failed" with
> error code 0x1A03. You have access to fleet metrics, server logs, and can
> SSH into test vehicles. Walk through your debugging process.

**Things to cover:**
- How do you triage? What do you look at first?
- How do you narrow down the 5%? What do they have in common?
- What tools would you use on a Linux system to diagnose?
- How do you distinguish between a server-side vs client-side issue?
- What's your communication plan while debugging?

<details>
<summary>Hints / discussion points</summary>

- **Triage**:
  1. Is the rollout still active? Pause it immediately if failure rate is above threshold.
  2. Look up error code 0x1A03 — what does it mean? (e.g., checksum mismatch,
     write failure, timeout)
  3. When did failures start? Correlate with rollout timeline.
- **Cohort analysis**:
  - Group failing vehicles by: hardware revision, current firmware version,
    geographic region, connectivity type (WiFi vs LTE).
  - Often the issue is version-specific (delta from version X is broken) or
    hardware-specific (different flash chip on older revisions).
- **On-device debugging**:
  - `dmesg` / `journalctl` for kernel and service logs.
  - Check disk space (`df -h`) — is the partition full?
  - `strace` the update process to see where it fails.
  - Verify the downloaded artifact: `sha256sum` against manifest.
  - Check mount points, filesystem integrity (`fsck`).
- **Server-side checks**:
  - Are the CDN artifacts correct? Download and verify checksums.
  - Are the delta patches valid for the source version these vehicles have?
  - Any server-side errors in the update API logs?
- **Communication**: Alert the rollout owner, post in incident channel, update
  status page. If safety-critical, escalate immediately.
</details>

---

## 9. Designing a Local Update Mechanism (USB/Ethernet)

> Some vehicles in the factory or service centers can't connect to the internet.
> Design a system that lets a technician update a vehicle via USB drive or
> local Ethernet connection.

**Things to cover:**
- How is the update package prepared and distributed to service centers?
- USB workflow: plug in drive → vehicle detects → applies update
- Ethernet workflow: technician laptop runs a local update server
- Security: how do you prevent unauthorized updates via USB?
- How does this integrate with the OTA system? (version tracking, reporting)

<details>
<summary>Hints / discussion points</summary>

- **Signed bundles**: Same update packages as OTA, signed by Tesla. The vehicle
  verifies the signature regardless of delivery mechanism. No special "USB key."
- **USB auto-detection**: `udev` rule triggers when a USB drive with a specific
  file layout is inserted. The update agent scans for valid, signed update
  bundles and presents them on the infotainment screen (or applies automatically
  in factory mode).
- **Factory mode vs service mode**:
  - Factory: fully automated, updates all components, no user interaction.
  - Service: technician selects which components to update, sees progress.
- **Local Ethernet server**: A lightweight Go HTTP server on the technician's
  laptop serves the update files. Vehicle's update agent points to a local URL
  instead of Tesla's CDN. Could use mDNS/Avahi for discovery.
- **Security**: USB updates must be signed. Optionally require a service token
  (time-limited, tied to the VIN) to prevent random USB drives from triggering
  updates. Factory mode requires physical access to a diagnostic port.
- **Reporting**: Vehicle stores update events locally. When it next connects
  to the internet, it reports the update to Tesla's fleet management system.
  Version tracking stays consistent regardless of update delivery method.
</details>

---

## 10. Rate Limiting and Prioritization of Update Downloads

> 2 million vehicles all try to download a 500MB update at the same time.
> How do you prevent overwhelming Tesla's infrastructure while ensuring
> critical security patches reach vehicles quickly?

**Things to cover:**
- Server-side rate limiting and scheduling
- Client-side backoff and jitter
- Priority tiers (security fix vs feature update vs map data)
- CDN architecture
- Vehicle-side download scheduling (WiFi vs cellular, time of day)

<details>
<summary>Hints / discussion points</summary>

- **Server-side scheduling**: Don't push to all vehicles at once. Server assigns
  each vehicle a random download window within a time range (e.g., over 48 hours).
  Vehicle checks in, server says "your download window is T+6h to T+8h."
- **Client-side jitter**: Even within a window, add random delay before first
  request. Prevents thundering herd.
- **Priority tiers**:
  - P0 (security/safety): Push immediately, override scheduling.
  - P1 (important feature): Normal staged rollout.
  - P2 (cosmetic/maps): Best-effort, can wait days.
- **CDN**: Use a global CDN (CloudFront/Akamai/Fastly) to serve update
  artifacts. Edge caching means vehicles download from nearby PoPs, not
  Tesla's origin servers.
- **Vehicle-side preferences**:
  - Only download large updates on WiFi (unless P0).
  - Prefer downloading at night when the car is parked and on WiFi.
  - Respect user settings ("download updates automatically" vs "ask first").
- **Bandwidth estimation**: Vehicle measures download speed. If on a slow
  connection, request smaller chunks or defer to a better connection.
- **Resume**: Support HTTP Range requests. A vehicle that starts downloading at
  home, drives to work, and reconnects should resume where it left off.
</details>

---

## Tips for Answering Design Questions

1. **Clarify scope first.** "Should I focus on the client side, server side, or
   both?" "Is this for a single vehicle or fleet-wide?" Interviewers love this.

2. **Start with the happy path.** Describe the normal flow before diving into
   failure handling. But **do** get to failure handling — that's where the
   interesting discussion is.

3. **Think about failure modes.** For every component, ask "what happens if
   this fails?" Power loss, network drop, corrupted download, full disk, etc.

4. **Mention trade-offs.** "We could do A or B. A is simpler but doesn't handle
   X. B is more complex but covers X. I'd start with A and add B if needed."

5. **Draw on the job description.** Mention CAN bus, SPI, code signing, fleet
   metrics — shows you read it and understand the domain.

6. **Use Go concurrency naturally.** When discussing orchestration, mention
   goroutines, channels, contexts, WaitGroups. This is your chance to show
   Go fluency in a systems context.

7. **It's an internship.** You're not expected to know everything. Saying
   "I'm not sure about the exact protocol, but I'd approach it by..." is
   much better than guessing confidently.
