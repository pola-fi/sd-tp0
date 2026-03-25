class DrawState:
    def __init__(self, expected_agencies=5):
        self.ready = False
        self.known_agencies = set()
        self.done_agencies = set()
        self.winners_by_agency = {}
        self.expected_agencies = expected_agencies if expected_agencies > 0 else 5

    def mark_agencies_seen(self, bets):
        for bet in bets:
            self.known_agencies.add(str(bet.agency))

    def register_done(self, agency_id):
        self.done_agencies.add(agency_id)
        return self._should_run_draw()

    def get_winners_if_ready(self, agency_id):
        if not self.ready:
            return False, []
        return True, list(self.winners_by_agency.get(agency_id, []))

    def publish_winners(self, winners_by_agency):
        self.winners_by_agency = winners_by_agency
        self.ready = True

    def _should_run_draw(self):
        return (not self.ready) and len(self.done_agencies) >= self.expected_agencies
