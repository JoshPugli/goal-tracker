import SwiftUI

struct HomeView: View {
    @StateObject private var vm = HabitViewModel()
    @State private var showSettings = false
    @State private var showTodaySheet = false
    @State private var dragOffset: CGFloat = 0

    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {

                // Counters from backend stats
                VStack(alignment: .leading, spacing: 12) {
                    CounterLine(value: vm.statsMonth?.completed ?? 0, total: vm.statsMonth?.total ?? 0, suffix: "m", size: 84, opacity: 1.0)
                    CounterLine(value: vm.statsWeek?.completed  ?? 0, total: vm.statsWeek?.total  ?? 0, suffix: "w", size: 72, opacity: 0.72)
                    CounterLine(value: vm.statsDay?.completed   ?? 0, total: vm.statsDay?.total   ?? 0, suffix: "d", size: 60, opacity: 0.5)
                }
                .padding(.horizontal)
                .padding(.top, 12)
                .padding(.leading, 24)

                Spacer(minLength: 0)

                // Bottom grabber that can be swiped up to open sheet
                VStack(spacing: 12) {
                    Capsule()
                        .fill(Color.secondary.opacity(0.4))
                        .frame(width: 40, height: 5)
                        .padding(.top, 8)

                    Button {
                        showTodaySheet = true
                    } label: {
                        HStack(spacing: 8) {
                            Image(systemName: "checkmark.circle.fill")
                            Text("\(vm.todayRemainingCount) remaining")
                                .font(.body)
                        }
                        .foregroundStyle(.secondary)
                        .padding(.vertical, 12)
                        .frame(maxWidth: .infinity)
                        .clipShape(RoundedRectangle(cornerRadius: 12))
                    }
                }
                .padding(.horizontal)
                .padding(.bottom, 24)
                .background(
                    Color(.systemBackground)
                        .ignoresSafeArea(edges: .bottom)
                        .shadow(color: .black.opacity(0.05), radius: 8, y: -2)
                )
                .gesture(
                    DragGesture(minimumDistance: 10, coordinateSpace: .local)
                        .onChanged { value in
                            dragOffset = value.translation.height
                        }
                        .onEnded { value in
                            if value.translation.height < -40 { // upward drag
                                showTodaySheet = true
                            }
                            dragOffset = 0
                        }
                )
            }
            .toolbar {
                // Center title/dropdown
                ToolbarItem(placement: .principal) {
                    Menu {
                        ForEach(vm.competitions) { comp in
                            Button {
                                vm.selectedCompetition = comp
                            } label: {
                                if comp == vm.selectedCompetition {
                                    Label(comp.name, systemImage: "checkmark")
                                } else {
                                    Text(comp.name)
                                }
                            }
                        }
                    } label: {
                        HStack(spacing: 4) {
                            Text(vm.selectedCompetition.name)
                                .font(.body)
                            Image(systemName: "chevron.down")
                                .font(.subheadline)
                        }
                        .foregroundStyle(.secondary)
                        .padding(.vertical, 6)
                        .padding(.horizontal, 10)
                        .clipShape(Capsule())
                    }
                }

                // Existing settings button
                ToolbarItem(placement: .topBarTrailing) {
                    Button {
                        showSettings = true
                    } label: {
                        Image(systemName: "gearshape")
                            .foregroundStyle(.secondary)
                    }
                    .accessibilityLabel("Settings")
                }
            }
        }
        .task { await vm.refreshAll() }
        .sheet(isPresented: $showSettings) {
            SettingsView()
                .presentationDetents([.medium])
        }
        .sheet(isPresented: $showTodaySheet) {
            TodayGoalsSheet(vm: vm)
                .presentationDetents([.medium, .large])
                .presentationDragIndicator(.visible)
        }
    }
}

// MARK: - Components

private struct CounterLine: View {
    let value: Int
    let total: Int
    let suffix: String
    let size: CGFloat
    let opacity: Double

    var body: some View {
        Text("\(value)/\(total)\(suffix)")
            .font(.system(size: size, weight: .bold, design: .rounded))
            .foregroundStyle(Color.primary.opacity(opacity))
            .frame(maxWidth: .infinity, alignment: .leading)
    }
}

private struct TodayGoalsSheet: View {
    @ObservedObject var vm: HabitViewModel

    var body: some View {
        NavigationStack {
            List {
                Section(header: Text("Today")) {
                    ForEach(vm.todayStates) { state in
                        HStack {
                            Text(state.goal.name)
                            Spacer()
                            Button {
                                Task { await vm.toggleToday(for: state.goal) }
                            } label: {
                                Image(systemName: state.completed ? "checkmark.circle.fill" : "circle")
                                    .foregroundStyle(state.completed ? .green : .secondary)
                                    .imageScale(.large)
                            }
                            .buttonStyle(.plain)
                        }
                    }
                }
            }
            .listStyle(.insetGrouped)
            .toolbar {
                ToolbarItem(placement: .principal) {
                    Text("\(vm.todayCompletedCount)/\(vm.habits.count) completed")
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                }
            }
        }
    }
}

// MARK: - ViewModel + Models

final class HabitViewModel: ObservableObject {
    // Backend models
    struct Goal: Identifiable, Codable, Hashable {
        let id: String
        let name: String
    }

    struct TodayState: Identifiable, Codable {
        var id: String { goal.id }
        let goal: Goal
        var completed: Bool
    }

    struct Stats: Codable {
        let window: String
        let completed: Int
        let total: Int
    }

    // Published state
    @Published var todayStates: [TodayState] = []
    @Published var statsDay: Stats?
    @Published var statsWeek: Stats?
    @Published var statsMonth: Stats?

    // Derived
    var habits: [Goal] { todayStates.map { $0.goal } }
    var todayCompletedCount: Int { todayStates.filter { $0.completed }.count }
    var todayRemainingCount: Int { max(habits.count - todayCompletedCount, 0) }

    // Networking
    // Development: point to ngrok tunnel for on-device testing
    private let baseURL = URL(string: "https://48458bf86bcf.ngrok-free.app")!

    func refreshAll() async {
        struct Dashboard: Codable {
            let stats_day: Stats
            let stats_week: Stats
            let stats_month: Stats
            let today: [TodayState]
        }
        do {
            let dash: Dashboard = try await request(path: "/api/dashboard")
            await MainActor.run {
                self.statsDay = dash.stats_day
                self.statsWeek = dash.stats_week
                self.statsMonth = dash.stats_month
                self.todayStates = dash.today
            }
        } catch {
            // Fallback to individual calls if dashboard not available
            async let day: Stats = fetchStats(window: "day")
            async let week: Stats = fetchStats(window: "week")
            async let month: Stats = fetchStats(window: "month")
            async let today: [TodayState] = request(path: "/api/goals/today")
            if let (d, w, m, t) = try? await (day, week, month, today) {
                await MainActor.run {
                    self.statsDay = d
                    self.statsWeek = w
                    self.statsMonth = m
                    self.todayStates = t
                }
            }
        }
    }

    func toggleToday(for goal: Goal) async {
        let isCompleted = todayStates.first(where: { $0.goal.id == goal.id })?.completed ?? false
        let method = isCompleted ? "DELETE" : "POST"
        let url = URL(string: "/api/goals/\(goal.id)/complete", relativeTo: baseURL)!
        var req = URLRequest(url: url)
        req.httpMethod = method
        do {
            // Optimistically update UI
            await MainActor.run {
                if let idx = self.todayStates.firstIndex(where: { $0.goal.id == goal.id }) {
                    self.todayStates[idx].completed.toggle()
                }
            }

            _ = try await URLSession.shared.data(for: req)

            // Pull fresh dashboard in one round trip
            struct Dashboard: Codable { let stats_day: Stats; let stats_week: Stats; let stats_month: Stats; let today: [TodayState] }
            let dash: Dashboard = try await request(path: "/api/dashboard")
            await MainActor.run {
                self.statsDay = dash.stats_day
                self.statsWeek = dash.stats_week
                self.statsMonth = dash.stats_month
                self.todayStates = dash.today
            }
        } catch {
            // Revert optimistic change on failure
            await MainActor.run {
                if let idx = self.todayStates.firstIndex(where: { $0.goal.id == goal.id }) {
                    self.todayStates[idx].completed.toggle()
                }
            }
            // TODO: handle error
        }
    }

    private func fetchStats(window: String) async throws -> Stats {
        let encoded = window.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? window
        return try await request(path: "/api/stats?window=\(encoded)")
    }

    private func request<T: Decodable>(path: String) async throws -> T {
        guard let url = URL(string: path, relativeTo: baseURL) else {
            throw URLError(.badURL)
        }
        var req = URLRequest(url: url)
        req.setValue("application/json", forHTTPHeaderField: "Accept")
        let (data, resp) = try await URLSession.shared.data(for: req)
        guard let http = resp as? HTTPURLResponse, (200..<300).contains(http.statusCode) else {
            throw URLError(.badServerResponse)
        }
        return try JSONDecoder().decode(T.self, from: data)
    }

    // Competition UI state (local for now)
    @Published var competitions: [Competition]
    @Published var selectedCompetition: Competition

    init() {
        let comps: [Competition] = [
            Competition(name: "Personal"),
            Competition(name: "Team A"),
            Competition(name: "Gym Buddies"),
        ]
        self.competitions = comps
        self.selectedCompetition = comps[0]
    }
}

// Removed old local models in favor of backend integration

struct Competition: Identifiable, Hashable {
    let id: UUID = UUID()
    var name: String
}

// MARK: - Preview

#Preview {
    HomeView()
}

import SwiftUI

struct SettingsView: View {
    @State private var remindersEnabled = true

    var body: some View {
        NavigationStack {
            Form {
                Toggle("Enable reminders", isOn: $remindersEnabled)
                NavigationLink("Manage habits") {
                    Text("Coming soon")
                        .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .center)
                }
            }
            .navigationTitle("Settings")
        }
    }
}

#Preview {
    SettingsView()
}

