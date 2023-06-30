## Inspiration

I have been following the Cosmos ecosystem for over two years now, and I would love to actively participate in its success.
However, one thing that always annoyed me was the lack of clarity about what is happening on-chain.
If I don't check my wallet regularly, I don't know if my validator got slashed, my unstaking period is over,
a new governance proposal is out, or if my borrow position got liquidated.
This led me to the idea of creating a tool that gives Cosmonauts personalized notifications about their on-chain activity.

## What it does

You can log in with your Keplr wallet into the [Star Scope webapp](https://starscope.network).
You are then notified about relevant events related to your wallet.

**Events** are triggered by an **action** of **you** or a **third party** (validator, another wallet, DAO, etc.).
Every **action** creates an **immediate** or **future event**.

|                     | Action by you | Action by third party |
|:-------------------:|:-------------:|:---------------------:|
| **Immediate event** |   (nothing)   |      Notify now       |
|  **Future event**   | Notify later  |     Notify later      |

### Examples
**Action by you, Immediate event**

This case is not covered since you can see the result of your action immediately.

**Action by you, Future event**
- You unstake your Osmo tokens. You get notified when the unstaking period is over.
- You unbond tokens in an Osmosis pool. You get notified when the unbonding period is over.

**Action by third party, Immediate event**
- Your validator gets slashed.
- A new governance proposal is out.
- Your borrow position gets liquidated.
- You receive tokens from someone else.
- Your sell order on a DEX (e.g. Injective) gets filled.

> You get notified immediately.

**Action by third party, Future event**
- Your validator falls out of the active set. You get notified if it doesn't get back in the active set after 48 hours.

## How I built it
I built most of the components in Golang because I am familiar with it, and the Cosmos SDK is written in Golang.
I used Rust for the frontend because I saw it as an opportunity to learn it.

### Tech stack

**Server:**
- A gRPC-web server written in Golang with a Postgres database. It is responsible for storing the data and serving it to the frontend.
- Event processors written in Golang that listen to a Kafka topic for new blockchain transactions and store them in the database.

**Indexers:**
- Osmosis indexer written in Golang that listens to new blocks on the Osmosis blockchain and publishes them to a Kafka topic.
- More indexers will follow (Injective, Mars, Neutron, etc.).

**Frontend:**
- A Rust application, that uses the Sycamore framework. For styling, I use TailwindCSS.

**Infrastructure:**
- A Docker Compose file that starts all the components and a Caddy server that serves the frontend.
- The full app is deployed on DigitalOcean.

**Architecture**

![Architecture](https://raw.githubusercontent.com/loomi-labs/star-scope/main/data/documentation/architecture.png)

## Challenges I ran into

If I use public blockchain nodes I have to deal with rate limits, slow responses and unreliable data.
I could solve this by running my own nodes but this is not feasible if I want to cover the whole Cosmos ecosystem.

## Accomplishments that I am proud of
I built a full prototype within a bit more than a week.

## What I learned
I learned a lot about how to build an indexer and connect everything together into a fullstack application.

## What's next for Star Scope
- Build a UI to show the events grouped by type.
- Add ability to get push notifications.
- More indexers to cover the whole Cosmos ecosystem.
- Add more notification channels (Telegram, Discord, etc.).