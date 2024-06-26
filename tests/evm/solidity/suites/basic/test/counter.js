/* eslint-disable no-undef */

const Counter = artifacts.require('Counter')

async function expectRevert(promise) {
  try {
    await promise
  } catch (error) {
    if (error.message.indexOf('revert') === -1) {
      expect('revert').to.equal(
        error.message,
        'Wrong kind of exception received'
      )
    }
    return
  }
  expect.fail('Expected an exception but none was received')
}

contract('Counter', (accounts) => {
  const [one, two, three] = accounts
  console.log('accounts', accounts)
  let counter

  beforeEach(async () => {
    const balanceOne = await web3.eth.getBalance(one)
    console.log('balanceOne', balanceOne)
    // send eth from account one to account two
    const receipt = await web3.eth.sendTransaction({
      from: one,
      to: two,
      value: web3.utils.toWei('1', 'ether')
    })
    console.log('receipt', receipt)
    await web3.eth.sendTransaction({
      from: one,
      to: three,
      value: web3.utils.toWei('1', 'ether')
    })

    counter = await Counter.new()
    console.log(counter.address)
  })

  it('should add', async () => {
    const balanceOne = await web3.eth.getBalance(one)
    const balanceTwo = await web3.eth.getBalance(two)
    const balanceThree = await web3.eth.getBalance(three)

    let count

    await counter.add({ from: one })
    count = await counter.getCounter()
    assert.equal(count, '1', 'Counter should be 1')
    assert.notEqual(
      balanceOne,
      await web3.eth.getBalance(one),
      `${one}'s balance should be different`
    )

    await counter.add({ from: two })
    count = await counter.getCounter()
    assert.equal(count, '2', 'Counter should be 2')
    assert.notEqual(
      balanceTwo,
      await web3.eth.getBalance(two),
      `${two}'s balance should be different`
    )

    await counter.add({ from: three })
    count = await counter.getCounter()
    assert.equal(count, '3', 'Counter should be 3')
    assert.notEqual(
      balanceThree,
      await web3.eth.getBalance(three),
      `${three}'s balance should be different`
    )
  })

  it('should subtract', async () => {
    let count

    await counter.add()
    count = await counter.getCounter()
    assert.equal(count, '1', 'Counter should be 1')

    // Use receipt to ensure logs are emitted
    const receipt = await counter.subtract()
    count = await counter.getCounter()

    assert.equal(count, '0', 'Counter should be 0')
    assert.equal(
      receipt.logs[0].event,
      'Changed',
      "Should have emitted 'Changed' event"
    )
    assert.equal(
      receipt.logs[0].args.counter,
      '0',
      "Should have emitted 'Changed' event with counter being 0"
    )

    // Check lifecycle of events
    const contract = new web3.eth.Contract(counter.abi, counter.address)
    const allEvents = await contract.getPastEvents('allEvents', {
      fromBlock: 1,
      toBlock: 'latest'
    })
    const changedEvents = await contract.getPastEvents('Changed', {
      fromBlock: 1,
      toBlock: 'latest'
    })
    assert.equal(allEvents.length, 3)
    assert.equal(changedEvents.length, 2)

    await expectRevert(counter.subtract())
  })
})
