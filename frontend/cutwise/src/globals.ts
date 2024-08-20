export const apiUrl: string = "http://localhost:8080/api/v1/";

export interface jobListJob {
  JobNumber: string;
  Customer: string;
}
export interface jobList {
  JobList: jobListJob[];
}
