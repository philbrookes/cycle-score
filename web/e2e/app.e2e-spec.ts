import { CycleScore.ComPage } from './app.po';

describe('cycle-score.com App', () => {
  let page: CycleScore.ComPage;

  beforeEach(() => {
    page = new CycleScore.ComPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
