import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MseComponent } from './mse.component';

describe('MseComponent', () => {
  let component: MseComponent;
  let fixture: ComponentFixture<MseComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MseComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MseComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
