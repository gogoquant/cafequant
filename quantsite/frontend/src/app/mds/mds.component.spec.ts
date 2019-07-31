import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MdsComponent } from './mds.component';

describe('MdsComponent', () => {
  let component: MdsComponent;
  let fixture: ComponentFixture<MdsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MdsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MdsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
