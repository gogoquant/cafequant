import { Component, OnInit } from '@angular/core';
@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {

  chartOption = {
    tooltip: {
      formatter: '{a} <br/>{b} : {c}%'
    },
    toolbox: {
      feature: {
        restore: {},
        saveAsImage: {}
      }
    },
    series: [
      {
        name: '业务指标',
        type: 'gauge',
        detail: { formatter: '{value}%' },
        data: [{ value: 50, name: '完成率' }]
      }
    ]
  };
  private echartsIntance: any;

  private runningLoop = () => {
    this.chartOption['series'][0]['data'][0]['value'] = parseFloat((Math.random() * 100).toFixed(2));
    console.log(this.chartOption['series'][0]['data'][0]['value']);
    // if (this.chartOption) {
    //   this.chartOption = this.chartOption;
    //   this.echartsIntance.setOption(this.chartOption, true);
    // }
  }
  constructor() { }

  ngOnInit() {
    this.runningLoop();
    // setInterval(this.runningLoop, 1000);
  }
  onChartInit(ec) {
    this.echartsIntance = ec;
  }
}
