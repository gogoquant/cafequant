import { AgileModule } from './agile.module';

describe('AgileModule', () => {
  let agileModule: AgileModule;

  beforeEach(() => {
    agileModule = new AgileModule();
  });

  it('should create an instance', () => {
    expect(agileModule).toBeTruthy();
  });
});
