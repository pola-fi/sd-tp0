import csv
import datetime
import threading


""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574

_STORAGE_LOCK = threading.Lock()


""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

""" Checks whether a bet won the prize or not. """
def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER

"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Process-safe guarantees are out of scope.
"""
def store_bets(bets: list[Bet]) -> None:
    with _STORAGE_LOCK:
        with open(STORAGE_FILEPATH, 'a+') as file:
            writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
            for bet in bets:
                writer.writerow([bet.agency, bet.first_name, bet.last_name,
                                 bet.document, bet.birthdate, bet.number])

"""
Load all bets from STORAGE_FILEPATH using a consistent snapshot.
Process-safe guarantees are out of scope.
"""
def load_bets() -> list[Bet]:
    with _STORAGE_LOCK:
        with open(STORAGE_FILEPATH, 'r') as file:
            reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
            rows = list(reader)

    return [Bet(row[0], row[1], row[2], row[3], row[4], row[5]) for row in rows]
