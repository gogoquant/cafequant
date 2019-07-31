import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TtaComponent } from './tta.component';

describe('TtaComponent', () => {
  let component: TtaComponent;
  let fixture: ComponentFixture<TtaComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TtaComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TtaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
