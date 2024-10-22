## Some go examples based on Domain-Driven Design with Golang
https://www.oreilly.com/library/view/domain-driven-design-with/9781804613450/

### 1) Recommendation Service

_Approach to build microservices with DDD principles._

*Context outline:*

### 2) CoffeeCo

_Approach to build monolithic-like application with DDD principles._

*Context outline:*

CoffeeCo is a national coffee shop chain. They experienced rapid growth in the last year and have opened 50 new stores. Each store sells coffee and coffee-related accessories, as well as store-specific drinks. Stores often have individual offers, but national marketing campaigns are often run, which influence the price of an item too.
CoffeeCo recently launched a loyalty program called CoffeeBux, which allows customers to get 1 free drink for every 10 they purchase. It doesn’t matter which store they purchase a drink at or which they redeem it at.

*Ubiquitous language:*
- Coffee lovers - what CoffeeCo calls its customers
- CoffeeBux - this is the name of their loyalty program. Coffee lovers earn one CoffeeBux for each drink or accessory they purchase
- Tiny, medium, and massive: The sizes of the drinks are in ascending order. Some drinks are only available in one size, others in all three. Everything on the menu fits into these categories

*Domains:*
- store
- purchase
- loyalty
- payment

*Features for MVP:*
- Purchasing a drink or accessory using CoffeeBux
- Purchasing a drink or accessory with a debit/credit card
- Purchasing a drink or accessory with cash
- Earning CoffeeBux on purchases
- Store-specific (but not national) discounts
- We can assume all purchases are in USD for now; in the future, we need to support many currencies though
- Drinks only need to come in one size for now
