# Module `x/distribution`

## **Overview**

The purpose of `distribution` module is responsible for distributing rewards between validators and delegators per `epoch` (`epoch` define on-chain timers, that execute fixed time interval). Additionally, the distribution module defines the community pool, which is a pool of funds under the control of on-chain governance.

The module in `aura` has 4 parameters that may be modified by governance proposal.

```
    1. communitytax: 0.02
    2. baseproposerreward: 0.04
    3. bonusproposerreward: 0.01
    4. withdrawaddrenabled: true
```

## **Content**

1. Concepts
2. State
3. Parameters
4. Query
5. Keeper
6. Hooks




### **1. Concepts**

Each validator has the opportunity to charge delegators commission on the rewards collected on behalf of the delegators. Fees are collected directly into a global reward pool and a validator proposer reward pool. Due to the nature of passive accounting, whenever changes to parameters which affect the rate of reward distribution occur, withdrawal of rewards must also occur.

**Formula**:

At epoch `e`, suppose there are `n` blocks with `v` validators (it means have `n` validator is proposer, `n < v`) and collects a total `T` in fees.

A validator participate in epoch has define by struct

```
ValidatorInfo {
    proposer_blocks: int
    active_blocks: int
    power: int
}
    
```

So validator `Vi` is participate `Mi = active_blocks - proposer_blocks` times on epoch as non-proposer

First a `communitytax` is applied. The fee go to the community pool (aka reserve pool). Reserve pool's funds can be allocated through governance to fund bouties and upgrades.

Let `R` is reward for each validator received in the epoch `e`.

We have:
```
n * [R + (baseproposerreward + bonusproposerreward) * R] + Sum(Mi) * R = T
```

$$=>R = {T \over [n + n*(baseproposerreward + bonusproposerreward) + Sum(Mi)]}$$

So

* For the proposer validator:

    The pool obtains `VR = R + R * (baseproposerreward + bonusproposerreward)`
    
    `commision = VR * (1-self_bonded) * commission_rate`
    
    Validator's reward: `VR * self_bonded + commission`

    Delegator's reward: `VR * (1-self_bonded) - commission`

* For the non-proposer validator:

    The pool obtains `NVR = R`
    
    `commision = NVR * (1-self_bonded) * commission_rate`
    
    Validator's reward: `NVR * self_bonded + commission`

    Delegator's reward: `NVR * (1-self_bonded) - commission`


Notes:

* All fees are distributed among all the bonded validators, in proportion to their consensus power.



