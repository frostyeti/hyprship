# Research: RDP Integration & Bastion Strategy

## Background
We have finalized the decision to use the native `golang.org/x/crypto/ssh` package for all SSH-related functionality within the API. This will power:
* Live SSH terminal sessions directly through the website/API.
* Job executions that run scripts on remote servers over SSH.
* Pushing and pulling files to/from remote servers via SCP.

However, we still need to determine the best approach for supporting live RDP (Remote Desktop Protocol) sessions in the browser. 

## Objectives
Investigate and compare the following solutions for bridging RDP connections to a web-based HTML5 canvas:

1. **Gravitational Teleport Integration:**
   * Research how easily Teleport can be integrated or deployed alongside our stack.
   * Can we use Teleport strictly for its RDP/Bastion capabilities while managing the UI/API orchestrations ourselves?
   * Evaluate the overhead and licensing considerations.

2. **Apache Guacamole (`guacd` sidecar):**
   * Evaluate running the `guacd` daemon as a sidecar/microservice.
   * Guacamole is the industry standard for translating RDP (and VNC/SSH) into a custom protocol that is easily rendered in an HTML5 canvas via WebSockets.
   * Assess the complexity of writing a Go client/bridge in our API that talks to `guacd`.

3. **Pure Go Implementation (`citilinkru/go-rdp`):**
   * Investigate the feasibility of using a pure Go RDP client library.
   * Determine if we would have to build our own RDP-to-WebSockets image rendering pipeline (which can be highly complex regarding bitmap caching, cursor updates, and keystrokes) or if there are existing wrappers we can leverage.

## Acceptance Criteria
* A brief comparison document (Pros/Cons, overhead, maintenance burden) for Teleport vs. Apache Guacamole vs. `go-rdp`.
* A recommended architecture for the RDP implementation.
* (Optional) A minimal Proof of Concept (PoC) demonstrating an RDP connection to a Windows machine rendered in a browser.