import { CicdModule } from './cicd.module';

describe('CicdModule', () => {
  let cicdModule: CicdModule;

  beforeEach(() => {
    cicdModule = new CicdModule();
  });

  it('should create an instance', () => {
    expect(cicdModule).toBeTruthy();
  });
});
