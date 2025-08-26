#include <iostream>
#include <vector>
#include <numeric>
#include <iomanip>
#include <random>

/*This project is a simple C++ program that demonstrates a classic Moving Average Crossover strategy for algorithmic 
trading. It is designed for educational purposes to show how a basic trading algorithm can be implemented, backtested
against simulated data, and executed.*/

// Calculates the Simple Moving Average (SMA) for a given period.
double calculateSMA(const std::vector<double>& prices, int period) {
    if (prices.size() < period) {
        return 0.0; // Not enough data to calculate SMA
    }
    // Sum the last 'period' number of prices
    double sum = std::accumulate(prices.end() - period, prices.end(), 0.0);
    return sum / period;
}

// Runs the moving average crossover trading simulation.
void runTradingSimulation(const std::vector<double>& marketData, int shortWindow, int longWindow) {
    std::cout << "--- Starting Algorithmic Trading Simulation ---\n"
              << "Strategy: Moving Average Crossover (" << shortWindow << "-day vs " << longWindow << "-day SMA)\n"
              << std::endl;

    std::vector<double> currentPrices;
    double portfolioValue = 10000.0; // Starting with $10,000 cash
    int sharesOwned = 0;
    bool positionOpen = false; // Are we currently holding shares?

    // To track the previous state of the moving averages
    double prevShortSMA = 0.0;
    double prevLongSMA = 0.0;

    std::cout << std::fixed << std::setprecision(2);

    for (size_t i = 0; i < marketData.size(); ++i) {
        currentPrices.push_back(marketData[i]);

        // Wait until we have enough data for the long window
        if (currentPrices.size() < longWindow) {
            continue;
        }

        double shortSMA = calculateSMA(currentPrices, shortWindow);
        double longSMA = calculateSMA(currentPrices, longWindow);

        // --- Crossover Logic ---
        // A "Golden Cross" (Buy Signal) occurs when the short-term SMA crosses ABOVE the long-term SMA.
        if (shortSMA > longSMA && prevShortSMA <= prevLongSMA && !positionOpen) {
            int sharesToBuy = portfolioValue / marketData[i];
            sharesOwned = sharesToBuy;
            portfolioValue -= sharesToBuy * marketData[i];
            positionOpen = true;
            std::cout << "Day " << i + 1 << " | Price: $" << marketData[i]
                      << " | BUY SIGNAL (Golden Cross)"
                      << " | Bought " << sharesOwned << " shares." << std::endl;
        }
        // A "Death Cross" (Sell Signal) occurs when the short-term SMA crosses BELOW the long-term SMA.
        else if (shortSMA < longSMA && prevShortSMA >= prevLongSMA && positionOpen) {
            double saleValue = sharesOwned * marketData[i];
            portfolioValue += saleValue;
            std::cout << "Day " << i + 1 << " | Price: $" << marketData[i]
                      << " | SELL SIGNAL (Death Cross)"
                      << " | Sold " << sharesOwned << " shares. Portfolio: $" << portfolioValue << std::endl;
            sharesOwned = 0;
            positionOpen = false;
        }

        // Update previous SMA values for the next iteration
        prevShortSMA = shortSMA;
        prevLongSMA = longSMA;
    }

    // At the end of the simulation, if we still hold a position, sell it at the last price.
    if (positionOpen) {
        double finalSaleValue = sharesOwned * marketData.back();
        portfolioValue += finalSaleValue;
        std::cout << "\nEnd of simulation. Selling remaining " << sharesOwned
                  << " shares at final price $" << marketData.back() << std::endl;
    }
    
    std::cout << "\n--- Simulation Complete ---" << std::endl;
    std::cout << "Final Portfolio Value: $" << portfolioValue << std::endl;
}

int main() {
    // --- Generate Simulated Market Data ---
    // This creates a series of random but somewhat continuous price movements.
    std::vector<double> marketData;
    std::random_device rd;
    std::mt19937 gen(rd());
    std::normal_distribution<> d(0.0, 1.5); // Small daily price changes
    double lastPrice = 100.0;

    for (int i = 0; i < 200; ++i) {
        double change = d(gen);
        lastPrice += change;
        if (lastPrice < 10.0) lastPrice = 10.0; // Prevent price from going too low
        marketData.push_back(lastPrice);
    }
    
    // Define the short and long windows for the moving averages
    int shortWindow = 10;
    int longWindow = 30;
    
    runTradingSimulation(marketData, shortWindow, longWindow);
    
    return 0;
}