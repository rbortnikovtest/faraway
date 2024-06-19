# Quotes with PoW rate limiter

This project consists of a simple REST server returning random quotes and a client solving the challenge and requesting quotes.

## Server
Server uses a challenge-response protocol for rate limiting of the incoming requests. If client doesn't provide the **_challenge_** and the **_proof_** for it in a /quote request, the server will return a new **_challenge_** and a **_difficulty_** back to the client.

The returned **_challenge_** string is a hex encoded sequence of bytes signed by the server's RSA key. To solve the challenge client has to find such a sequence of bytes so that the SHA256 hash of the concatenation of the **_challenge_** and the **_proof_** **_SHA256(challenge+proof)_** would have N=**_difficulty_** first bytes equal to 0.

## Client
Client is a simple Go script which requests /quote from the server in an infinite loop. If the server returns a new **_challenge_** and a **_difficulty_** the client will solve the challenge and send the **_proof_** back to the server in the next /quote request. All the steps are logged to the console for better visibility.
