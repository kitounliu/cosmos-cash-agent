## Cosmos Cash controllers

The [Cosmos Cash Protocol](https://github.com/allinbits/cosmos-cash) is a Decentralized Key Management System (DKMS) which needs agents to implement certain protocols to allow the project to leverage Self sovereign identity [SSI](https://en.wikipedia.org/wiki/Self-sovereign_identity)

### Key Characteristics

- A [blockchain interface layer](https://github.com/allinbits/cosmos-cash-resolver) (known as a resolver) for creating and signing blockchain transactions.
- A resolver can be seen as part of a larger component known as the [VDR](https://github.com/allinbits/cosmos-cash):
  > Aries Verifiable Data Registry Interface: An interface for verifying data against an underlying ledger.
- A cryptographic wallet that can be used for secure storage of cryptographic secrets and other information (the secure storage tech, not a UI) used to build blockchain clients.
- An encrypted messaging system for allowing off-ledger interaction between those clients using multiple transport protocols.
- An implementation of ZKP-capable W3C verifiable credentials using the ZKP primitives found in Ursa.
- An implementation of the Decentralized Key Management System (DKMS) specification currently being incubated in under the name Cosmos Cash by AllinBits.
- A mechanism to build higher-level protocols and API-like use cases based on the secure messaging functionality.

#### Protocols

The following protocols are needed in the Cosmos Cash project.

#### 1. DIDExchange Protocol

#### 2. Introduce Protocol

#### 3. IssueCredential Protocol

#### 4. KMS

#### 5. Mediator

#### 6. Messaging

#### 7. OutOfBand Protocol

#### 8. PresentProof Protocol

#### 9. VDR

#### 10. Verifiable

### How to run

- `./scripts/start-agent.sh`