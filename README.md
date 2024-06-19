# Quotes with PoW rate limiter

This project provides a simple REST server returning random quotes.

It uses a challenge-response protocol for rate limiting the incoming requests. If client doesn't provide the **_challenge_** and the **_proof_** for it in a /quote request, the server will return a new **_challenge_** and a **_difficulty_** back to the client.

The returned **_challenge_** string is a hex encoded sequence of bytes signed by the server's RSA key. To solve the challenge client has to find such a sequence of bytes so that the SHA256 hash of the concatenation of the **_challenge_** and the **_proof_** (SHA256(challenge+proof)) would have N=**_difficulty_** first bytes equal to 0.
