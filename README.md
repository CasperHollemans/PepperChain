![PepperChain logo](logo.png "PepperChain Logo")

# PepperChain

PepperChain 🌶️, is a exploratory blockchain project to learn and get hands-on experience. It's focused on the
implementing the fundamental aspects of a blockchain.

## Features

This project is a simplistic implementation of a blockchain network. It's not intended to be used in production. The
features are chosen due to their importance in a blockchain network. The primary focus is to be able to make
transactions in a multi-node network.

### Blockchain

The blockchain data structure consists of a series of blocks in linear order. Each block contains a hash of the previous
block. This creates a chain of blocks, hence the name blockchain. Each block contains data, this can be transactions or
other types of data. A new block is created in a interval that can vary per network and the current state of the
network.
When a new block is created it is broadcast to the network, ensuring the blockchain stays up to date in each node.

A blockchain is a immutable structure, each block that is added strengthens the immutability of all previous blocks.
Every block contains the hash of a previous block, if one block in the chain is changed the hash will change as well.
This means that all subsequent block have to be updated in order to maintain a valid chain.

### Proof of Work

Proof of work is a consensus algorithm used to verify transactions and create new blocks. When a new block is created,
the node is supposed to mine the block. This will be done by solving a complex puzzle or challenge. The goal of the
challenge is to generate a hash that meets the current difficulty of the network. Proof of Work is considered a secure
consensus mechanism in well established blockchain networks. Due to the complexity and randomness of the puzzle it
takes at least 51% of the complete computational power of a network to gain control and manipulate the blockchain.

### Transactions

In blockchain a transaction represents the movement of a digital asset from one wallet to another. By tracing all
transactions that are associated with your wallet address you can determine your current balance or property.

When a new transaction is created it will be verified by the network. The network will verify if the transaction is
valid by checking the digital signature of the transaction and balance of the sender. After verification the node(s)
that verified the transaction will propagate it further to other nodes for verification. Finally the transaction is added
to the memory pool, waiting to be included in a block.

### Peer to Peer Network

Having a well stabilised network is a crucial part of the security in a blockchain. Having a decentralised network
enhances the security because there is no single point of failure. This makes directed attacks at single nodes obsolete
because the rest of the network has to agree to the changes that are made to the blockchain. Like mentioned earlier,
you'd need at least 51% to gain control over the blockchain.

A new node will register themselves in the network, making known that they exist and how to reach them. After this the
network will start communicating with the node. In order to communicate back to the network the node needs to download
a list of all other network nodes.

When communicating between peers in the network the node randomly selects some other peers in the network and send the
relevant information. The receiving node will do the same creating a chain reaction, eventually all nodes will be up to
date. This is called the gossip protocol, a way of sharing information in a distributed system.

## Acknowledgements

- [Blockchain in Go By Jeiwan](https://github.com/Jeiwan/blockchain_go)