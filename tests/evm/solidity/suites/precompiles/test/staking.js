const { expect } = require('chai')
const hre = require('hardhat')

// const E12 = hre.ethers.BigNumber.from('1000000000000')
const stakingAddr = '0x0000000000000000000000000000000000000800'
const valAddr = 'auravaloper1dd6psq88kntuzyap944et8fmh0mxmw2wqnrnpx'

describe('Staking', function () {
  it('should stake to a validator', async function () {
    const stakeAmount = hre.ethers.parseEther('0.000111')

    const [signer] = await hre.ethers.getSigners()
    const staking = await hre.ethers.getContractAt(
      'StakingI',
      stakingAddr,
      signer
    )
    console.log('Stake amount:', stakeAmount.toString())


    // get original balance
    const originalBalance = await hre.ethers.provider.getBalance(signer)
    console.log('Original balance:', originalBalance)

    // get current delegation
    const currentDelegation = await staking.delegation(signer.address, valAddr)
    console.log('Current delegation:', currentDelegation)

    const tx = await staking
      .delegate(signer.address, valAddr, stakeAmount)
    const res = await tx.wait()

    // Query delegation
    const delegation = await staking.delegation(signer.address, valAddr)
    console.log('Delegation:', delegation)

    expect(delegation.balance.amount).to.equal(
      stakeAmount + currentDelegation.balance.amount,
      'Stake amount does not match'
    )

    // check balance
    const newBalance = await hre.ethers.provider.getBalance(signer)
    console.log('New balance:', newBalance);
    // new balance should be less than original balance - stake amount (because of gas fee)
    expect(newBalance).to.lessThan(originalBalance - stakeAmount, 'Available amount does not match')
  })

  it('should undelegate from a validator', async function () {
    const [signer] = await hre.ethers.getSigners()
    const staking = await hre.ethers.getContractAt(
      'StakingI',
      stakingAddr,
      signer
    )

    const prevDelegation = await staking.delegation(signer.address, valAddr)
    console.log('Prev delegation:', prevDelegation)

    // undelegate half of the amount
    const tmp = BigInt(prevDelegation.balance.amount) / 2n
    const undelegateAmount = tmp - (tmp % (10n ** 12n))
    console.log('Undelegate amount:', undelegateAmount)

    const tx = await staking.undelegate(
      signer.address,
      valAddr,
      undelegateAmount,
      // TODO: gas estimation is not working, we have to set gas limit manually
      { gasLimit: 200000 }
    )
    const res = await tx.wait()

    const delegation = await staking.delegation(signer.address, valAddr)
    console.log('Delegation:', delegation)

    expect(delegation.balance.amount).to.equal(
      prevDelegation.balance.amount - undelegateAmount,
      'Undelegated amount does not match'
    )

    const undelegation = await staking.unbondingDelegation(signer.address, valAddr)
    console.log('Unbonding delegation:', undelegation)

    expect(undelegation.balance.amount).to.equal(
      undelegateAmount,
      'Unbonding amount does not match'
    )
  })
})
