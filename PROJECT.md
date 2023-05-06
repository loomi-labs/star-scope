## Inspiration

I follow the Cosmos ecosystem since over 2 years now and I would love to actively participate in its success.
Something that always annoyed me was the lack of clarity of what is happening on-chain. If I don't check my wallet
regularly I don't know if my validator got slashed, my unstaking period is over, a new governance proposal is out or if
my borrow position got liquidated. This led me to the idea of creating a tool that gives Cosmonauts personalized
notifications about their on-chain activity.

## What it does

You can login with your Keplr wallet into the Star Scope webapp. You get then notified about relevant events related to
your wallet.

**Events** are triggered by an **action** of **you** or a **3-party** (validator, another wallet, DAO, etc.).\
Every **action** creates an **immediate** or **future event**.

|                     | Action by you | Action by 3-party |
|---------------------|---------------|-------------------|
| **Immediate event** | (nothing)     | Notify now        |
| **Future event**    | Notify later  | Notify later      |

### Examples
**Action by you, Immediate event**\
This case is not covered since you can see the result of your action immediately.

**Action by you, Future event**
- You unstake your Osmo tokens. You get notified when the unstaking period is over.
- You unbond tokens in a Osmosis pool. You get notified when the unbonding period is over.

**Action by 3-party, Immediate event**
- You validator gets slashed. You get notified immediately.
- A new governance proposal is out. You get notified immediately.
- You borrow position gets liquidated. You get notified immediately.
- You receive tokens from someone else. You get notified immediately.
- Your sell order on a DEX (e.g. Injective) gets filled. You get notified immediately.

> You get notified immediately.

**Action by 3-party, Future event**
- Your validator falls out of the active set. You get notified if he doesn't get back in the active set after 48 hours.

## How I built it
I built most ot the components in golang because I am familiar with it and Cosmos SDK is written in golang.
I used Rust for the frontend because I saw it as a good opportunity to learn it.

### Tech stack

**Server:**
- A gRPC-web server written in golang with a postgres database. It is responsible for storing the data and serving it to the frontend.
- Event processors written in golang that listen to a kafka topic (`indexed-events`) for new 
blockchain transactions and stores them in the database or publishe them to a kafka topic (`processed-events`.

**Indexers:**
- Osmosis indexer written in golang that listens to new blocks on the osmosis blockchain and publishes them to a kafka topic (`indexed-events`).
- More indexers will follow (Injective, Mars, Neutron, etc.)

**Frontend:**
- A rust application that uses the [Sycamore](https://sycamore-rs.netlify.app/) framework. For styling I use [tailwindcss](https://tailwindcss.com/).

**Infrastructure:**
- A docker-compose file that starts all the components and a caddy server that serves the frontend.
- The full app is deployed on DigitalOcean.

**Architecture**

![Architecture](https://raw.githubusercontent.com/loomi-labs/star-scope/233ebaa67f4299b2a8f86ab78fb8f09e1736f83c/data/documentation/architecture.png?token=GHSAT0AAAAAAB5XCFR2XMQNL3TUPSVWUO6SZCWFMMA)

## Challenges I ran into

If I use public blockchain nodes I have to deal with rate limits and unreliable data.\
I could solve this by running my own nodes but this is not feasible if I want to cover the whole Cosmos ecosystem.

## Accomplishments that I am proud of
I built a full prototype within a bit more than a week. 

## What we learned


## What's next for Star Scope
