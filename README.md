# GoBackTester
This is a backtester written in Golang (and a little python) that uses forex data in order to backtest given trading strategies. Our Data is sourced from truefx.com's historical forex data and supported currencies are listed on their website.
> Update Jan 1, 2023: Had to use some python and selenium to parse their site. Must integrate cmd line arguments and make loadcsv use this new file. Also must handle deleting the files directory each run and make waiting for page to load more efficient. It's 12:08AM though so happy new year I'm going to bed.
