import { MseModule } from './mse.module';

describe('MseModule', () => {
  let mseModule: MseModule;

  beforeEach(() => {
    mseModule = new MseModule();
  });

  it('should create an instance', () => {
    expect(mseModule).toBeTruthy();
  });
});
