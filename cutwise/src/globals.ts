export const apiUrl: string = "http://localhost:2828/api/v1/";

export interface jobListJob {
  JobNumber: string;
  Customer: string;
  Star: number;
}
export interface jobList {
  JobList: jobListJob[];
}


