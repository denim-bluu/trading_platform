# Trading System

## Service Responsibilities

1. ✅ Data Service:

   - Fetch historical stock data from external sources (e.g., Yahoo Finance API)
      - Store the response into Postgres database for future retrieval
      - Store the response into Cache (In-memory) for faster retrieval
   - Provide both single stock and batch stock data retrieval
   - Return OHLC (Open, High, Low, Close) price data, adjusted close prices, and volume
   - Handle date range and interval specifications in requests

2. ✅ Momentum Strategy Service:

   - Calculate momentum scores using exponential regression and R-squared for 90-day periods
   - Generate buy/hold signals based on momentum and price relative to 100-day moving average
   - Calculate ATR (Average True Range) for position sizing
   - Implement 15% gap check for the past 90 days
   - Rank stocks based on momentum scores
   - Filter to keep only the top 20% of stocks
   - Calculate position sizes based on ATR for risk parity
   - Provide a batch processing capability for multiple stocks
   - Sort and return signals for the top-ranked stocks
   - Logging of strategy calculations and decisions

3. Portfolio Service:
   - Manages the current portfolio composition
   - Handles position sizing and risk management
   - Performs weekly portfolio rebalancing
   - Tracks performance and generates reports

4. TODO: Trade Execution Service:
   - Interfaces with various brokers and exchanges
   - Executes trades based on signals from the Portfolio Service
   - Handles order management and trade confirmation
   - Provides real-time trade status updates

5. TODO: Backtesting Service:
   - Simulates trading strategies on historical data
   - Generates performance reports and statistics

6. TODO: API Gateway:
   - Provides a unified entry point for external requests

7. TODO: Scheduler Service:
   - Manages the timing of various trading activities
   - Triggers weekly trading actions (every Wednesday)
   - Initiates bi-weekly position size rebalancing
   - Schedules regular data updates and system maintenance tasks

```mermaid
graph TD
    A[Scheduler Service] -->|Triggers weekly| B[Portfolio Service]
    A -->|Triggers bi-weekly| B
    B -->|Requests data| C[Data Service]
    B -->|Requests signals| D[Strategy Service]
    B -->|Sends trade orders| E[Trade Execution Service]
    C -->|Provides market data| B
    C -->|Provides market data| D
    D -->|Provides signals| B
    E -->|Executes trades| F[Broker/Exchange]
```

## Trading Process Summary

1. Every Wednesday: Update Portfolio
   - Check S&P 500 relative to its 200MA
   - For existing positions:
     - Sell if:
       - No longer in top 20% of momentum-ranked stocks
       - Fallen below its 100MA
       - Received an explicit sell signal
     - Adjust size if:
       - Still in top 20% and above 100MA, but target position size has changed
   - For new positions:
     - Only open if S&P 500 is above its 200MA
     - Buy stocks from the top 20% that aren't already in the portfolio

2. Every Second Wednesday of the Month: Rebalance Portfolio
   - Perform all actions from the weekly update
   - Additionally:
     - Rebalance all existing positions to their target sizes based on current risk factors
     - This rebalancing occurs regardless of the S&P 500's position relative to its 200MA
