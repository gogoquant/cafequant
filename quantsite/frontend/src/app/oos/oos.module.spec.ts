import { OosModule } from './oos.module';

describe('OosModule', () => {
  let oosModule: OosModule;

  beforeEach(() => {
    oosModule = new OosModule();
  });

  it('should create an instance', () => {
    expect(oosModule).toBeTruthy();
  });
});
