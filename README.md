# EOS Light Node

[中文版](./README_cn.md)

## 0. Introduction

`EOS Light Node` is a [eosio](https://github.com/EOSIO/eos) node implementation develop by golang. Unlike the full node of eosio based on nodeos, the node by eos-light-node is only Verify the block header of the eosio chain.The actions in block will not be executed, so the light nodes have low requirements.

## 1. Motivation

EOSIO's block header structure supports lightweight nodes, but currently there is no reliable lightweight node implementation in the EOSIO community. In the EOSIO network, a full node has very high requirements for machine performance, network,  and operation  maintenance. At present, there are very few full nodes available in the EOSIO network. This means that the current EOSIO network structure is strongly centralized. It is has a big impact to stability, and it also makes developers harder to develop and deploy DAPP on the chain.

Although EOSIO's plugin architecture can allow nodes to expand functions, but nodes have higher requirements for machine performance, at the same time, developing plugins based on C ++ requires high development thresholds. Therefore, EOSIO has few expansion tools, which indirectly make development of DAPP more Difficulty.

At present, an important direction of chain development is to support cross-chain. EOSIO supports adding side chains based on IBC. IBC needs to be based on full nodes so that users can obtain the proof data which based on the merkle tree. However, due to the cost of a full node is too high, even if it is unrealistic to require a small number of validators and phishers to start a full node. This will make some theoretical support schemes can not work well as it is uneconomical. If all services only provided by full nodes, then the entire system will be very vulnerable.

To sum up the above problems, we need a sufficiently lightweight EOSIO node implement.Light nodes should be easy to expand and secondary development and ensure the verification of the blocks.

## 2. Design

One of the main reasons for EOSIO's high demand for nodes is its memory (RAM) status. These states are very large and have not been verified in the block, so any node wants to obtain any state in RAM, it needs to completely replay the EOSIO block . Although the node can be restored through the snapshot function, even for the nodes that are kept in sync, the recovery time is relatively long if there is a failure. At the same time, due to the busy transactions on the EOSIO chain, the peak value can reach about 4k+ TPS, and the synchronous real-time EOSIO block will also have higher requirements on the machine.

Based on the above considerations, EOS light nodes will not process every transaction, but only verify the validity of blocks and block headers, and confirm the blocks are irreversible.

Not processing every transactions means most of the work of light nodes can be completed in parallel, and synchronization can also be started from any block height. In the current prototype implementation, light nodes only need to use very few computing resources to complete verification.

Although not processing each transactions allows light nodes to occupy few resources, it also makes light nodes unable to provide EOS memory (RAM) state. In EOS, almost all states are in memory. Some operations (EOS calls it Action) also needs to be triggered by executing a transaction. This require some assertion mechanism to ensure that the node can reliably obtain the state of EOS based on the block data only. This can be achieved through the EOS state assertion contract and the state static assertion contract.

According to the different roles of EOS light nodes, we divide EOS light nodes into block nodes and block head nodes. The former includes all or part of EOS blocks, and the latter includes only EOS block header data.

Considering that EOSIO is still under development, EOS light nodes also need to follow it, so the design of EOS light nodes should be as homogeneous as possible with eosio, and its APIs and operations also should be as uniform as possible with nodeos, of course , Light nodes do not contain memory state, so some APIs, such as `get table`, can not provide by eos light node.

Light nodes are developed by golang, which can make light nodes easily run in different environments. At the same time, golang has many applications in back-end development, and is also suitable for applications such as developing data services based on light nodes.

## 3. Getting Start

EOSIO light nodes are developed by golang. First, you can install the golang environment by referring to [golang](https://golang.org/dl/)

### 3.1 Compile

```bash
git clone https://github.com/eosforce/eos-light-node.git
cd eos-light-node
go build
```

### 3.2 Start

Here assumed that the available p2p url of eosio is `127.0.0.1:9001` and the chain-id is `1c6ae7719a2a3b4ecb19584a30ff510ba1b6ded86e1fd8b8fc22f1179c622a32`

```bash
./eos-light-node -v -genesis "./eosio/config/genesis.json" -chain-id "1c6ae7719a2a3b4ecb19584a30ff510ba1b6ded86e1fd8b8fc22f1179c622a32" -p2p "127.0.0.1:9001"
```

> *Note*: Since there are few eosio p2p addresses available, you can start a local testnet based on the eosio boot script to test.

## 4. TODOs

Currently EOSIO light nodes are still under development, and the following functions will be completed in the future:

- Improve block storage implementation, compatible with nodeos' block db
- Implement irreversible block determination algorithm consistent with nodos
- Support other nodes to synchronize blocks and transactions from light nodes
- Improve API interface to be as compatible as possible with nodeos
- Optimize concurrent computing capabilities and speed up the synchronization process
- Support eos-vm for the implementation of some actions
