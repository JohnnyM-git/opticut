import {apiUrl} from "../globals.ts";

export async function getLocalJobs() {
    const res = await fetch(`${apiUrl}localjobs`)
    // console.log(await res.json())
    return await res.json()

}