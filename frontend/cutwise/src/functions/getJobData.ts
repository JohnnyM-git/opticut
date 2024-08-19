import {apiUrl} from "../globals.ts";

export async function getJobData(job :string):Promise<any> {
    const res = await fetch(`${apiUrl}job?job_id=${job}`)
    console.log(await res.json())
    // return await res.json()

}