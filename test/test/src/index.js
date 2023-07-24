const { Serverless } = require('universecloud-js')
const { Api } = require('@universecloud/api')
const { AuthType } = require('@universecloud/api/dist/api/client')

const wait = (ms) => {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve(true)
    }, ms)
  })
}
Serverless.handler(async function(invokeParams) {
  console.log(`hello from ${process.env.FUNCTION_ID}`)
  return {
    status: 204,
    body: ''
  }
})