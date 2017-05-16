import { Component } from '@angular/core';
import { ApiService } from './api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: []
})

export class AppComponent {
  public access_url;
  public authorized;
  public getting_score;
  public score;
  public startDate;
  public endDate;
  public error;
  constructor(private apiService: ApiService){}

  public resetForm() {
    this.getting_score = false;
    this.score = false;
    this.error = false;
  }

  public printScore():string {
    return Number(this.score.score).toLocaleString();
  }

  public getScore() {
    this.getting_score = true;
    let from = Date.parse(this.startDate) / 1000;
    let to = Date.parse(this.endDate) / 1000;
    this.apiService.getScore(from, to)
      .then((response) => {
        console.log(response);
        this.score = response;
        this.getting_score = false;
      })
      .catch((err) => {
        this.error = err;
        this.getting_score = false;
      });
  }

  ngOnInit() {
    this.apiService.isAuthorized()
      .then((response) => {
        if(response.valid === false) {
          this.apiService.oauthUrl()
          .then((response) => {
            this.access_url = response.url;
          })
          .catch((err) => {
            console.log("Error retrieving oauth URL:");
            console.log(err);
          });
        } else {
          this.authorized = true;
        }
      });

  }
}
