import { login } from './lib.mjs';
const [email, tag] = process.argv.slice(2);
const { b } = await login(email, tag); await b.close(); console.log('state', tag);
