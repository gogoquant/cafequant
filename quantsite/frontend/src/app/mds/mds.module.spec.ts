import { MdsModule } from './mds.module';

describe('MdsModule', () => {
  let mdsModule: MdsModule;

  beforeEach(() => {
    mdsModule = new MdsModule();
  });

  it('should create an instance', () => {
    expect(mdsModule).toBeTruthy();
  });
});
