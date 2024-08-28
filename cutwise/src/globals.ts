// globals.ts
export const apiUrl: string = "http://localhost:2828/api/v1/";

export interface JobListJob {
  JobNumber: string;
  Customer: string;
  Star: number;
}

export interface JobList {
  JobList: JobListJob[];
}

export interface Settings {
  kerf: number;
  units: string;
}
