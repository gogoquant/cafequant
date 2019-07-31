export class PodInfo {
  name = '';
  image = '';
  start_time = '';
  status = '';

  constructor(name: string, image: string, status: string, start_time: string) {
    this.name = name;
    this.image = image;
    this.start_time = start_time;
    this.status = status;
  }
}
