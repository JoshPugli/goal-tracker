import SwiftUI

final class HabitViewModel: ObservableObject {
    struct Goal: Identifiable, Codable, Hashable { let id: String; let name: String }
    struct TodayState: Identifiable, Codable { var id: String { goal.id }; let goal: Goal; var completed: Bool }
    struct Stats: Codable { let window: String; let completed: Int; let total: Int }

    @Published var todayStates: [TodayState] = []
    @Published var statsDay: Stats?
    @Published var statsWeek: Stats?
    @Published var statsMonth: Stats?

    var habits: [Goal] { todayStates.map { $0.goal } }
    var todayCompletedCount: Int { todayStates.filter { $0.completed }.count }
    var todayRemainingCount: Int { max(habits.count - todayCompletedCount, 0) }

    private let baseURL = AuthManager.shared.baseURL

    func refreshAll() async {
        struct Dashboard: Codable { let stats_day: Stats; let stats_week: Stats; let stats_month: Stats; let today: [TodayState] }
        do {
            let dash: Dashboard = try await request(path: "/api/dashboard")
            await MainActor.run {
                self.statsDay = dash.stats_day
                self.statsWeek = dash.stats_week
                self.statsMonth = dash.stats_month
                self.todayStates = dash.today
            }
        } catch {
            async let day: Stats = fetchStats(window: "day")
            async let week: Stats = fetchStats(window: "week")
            async let month: Stats = fetchStats(window: "month")
            async let today: [TodayState] = request(path: "/api/goals/today")
            if let (d, w, m, t) = try? await (day, week, month, today) {
                await MainActor.run { self.statsDay = d; self.statsWeek = w; self.statsMonth = m; self.todayStates = t }
            }
        }
    }

    func toggleToday(for goal: Goal) async {
        let isCompleted = todayStates.first(where: { $0.goal.id == goal.id })?.completed ?? false
        let method = isCompleted ? "DELETE" : "POST"
        let url = URL(string: "/api/goals/\(goal.id)/complete", relativeTo: baseURL)!
        var req = URLRequest(url: url)
        AuthManager.shared.attachAuth(to: &req)
        req.httpMethod = method
        do {
            await MainActor.run { if let idx = self.todayStates.firstIndex(where: { $0.goal.id == goal.id }) { self.todayStates[idx].completed.toggle() } }
            _ = try await URLSession.shared.data(for: req)
            struct Dashboard: Codable { let stats_day: Stats; let stats_week: Stats; let stats_month: Stats; let today: [TodayState] }
            let dash: Dashboard = try await request(path: "/api/dashboard")
            await MainActor.run { self.statsDay = dash.stats_day; self.statsWeek = dash.stats_week; self.statsMonth = dash.stats_month; self.todayStates = dash.today }
        } catch {
            await MainActor.run { if let idx = self.todayStates.firstIndex(where: { $0.goal.id == goal.id }) { self.todayStates[idx].completed.toggle() } }
        }
    }

    private func fetchStats(window: String) async throws -> Stats {
        let encoded = window.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? window
        return try await request(path: "/api/stats?window=\(encoded)")
    }

    private func request<T: Decodable>(path: String) async throws -> T {
        guard let url = URL(string: path, relativeTo: baseURL) else { throw URLError(.badURL) }
        var req = URLRequest(url: url)
        AuthManager.shared.attachAuth(to: &req)
        req.setValue("application/json", forHTTPHeaderField: "Accept")
        let (data, resp) = try await URLSession.shared.data(for: req)
        guard let http = resp as? HTTPURLResponse, (200..<300).contains(http.statusCode) else { throw URLError(.badServerResponse) }
        return try JSONDecoder().decode(T.self, from: data)
    }

    @Published var competitions: [Competition]
    @Published var selectedCompetition: Competition
    init() {
        let comps: [Competition] = [Competition(name: "Personal"), Competition(name: "Team A"), Competition(name: "Gym Buddies")]
        self.competitions = comps
        self.selectedCompetition = comps[0]
    }
}

struct Competition: Identifiable, Hashable { let id: UUID = UUID(); var name: String }


