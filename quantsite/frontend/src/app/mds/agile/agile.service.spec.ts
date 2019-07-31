import { TestBed, inject } from '@angular/core/testing';

import { AgileService } from './agile.service';

describe('AgileService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AgileService]
    });
  });

  it('should be created', inject([AgileService], (service: AgileService) => {
    expect(service).toBeTruthy();
  }));
});
