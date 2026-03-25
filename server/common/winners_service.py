from .utils import has_won, load_bets


class WinnersService:
    def compute_winners_by_agency(self):
        winners_by_agency = {}
        for bet in load_bets():
            if has_won(bet):
                agency_id = str(bet.agency)
                winners_by_agency.setdefault(agency_id, []).append(str(bet.document))
        return winners_by_agency

