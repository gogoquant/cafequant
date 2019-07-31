import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppInfoDetailComponent } from './app-info-detail.component';

describe('AppInfoDetailComponent', () => {
  let component: AppInfoDetailComponent;
  let fixture: ComponentFixture<AppInfoDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppInfoDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppInfoDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
