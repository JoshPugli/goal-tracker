import SwiftUI

struct HomeView: View {
    @StateObject private var vm = HabitViewModel()
    @State private var showSettings = false
    @State private var showTodaySheet = false
    @State private var dragOffset: CGFloat = 0

    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {

                // Counters: X/total with bigger type and padding
                VStack(alignment: .leading, spacing: 12) {
                    CounterLine(value: vm.completedCount(.month), total: vm.totalCapacity(.month), suffix: "m", size: 84, opacity: 1.0)
                    CounterLine(value: vm.completedCount(.week),  total: vm.totalCapacity(.week),  suffix: "w", size: 72, opacity: 0.72)
                    CounterLine(value: vm.completedCount(.day),   total: vm.totalCapacity(.day),   suffix: "d", size: 60, opacity: 0.5)
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
                            Text(state.habit.name)
                            Spacer()
                            Button {
                                vm.toggleToday(for: state.habit)
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
    enum Window { case day, week, month }

    @Published var habits: [Habit] = [
        Habit(name: "Drink water"),
        Habit(name: "Read 20 min"),
        Habit(name: "Exercise"),
    ]

    // All completion events across time
    @Published private(set) var completions: [HabitCompletion] = SampleData.seededCompletions

    // Derived: today’s state for each habit
    var todayStates: [HabitState] {
        habits.map { habit in
            HabitState(habit: habit, completed: isCompletedToday(habit))
        }
    }

    var todayCompletedCount: Int {
        todayStates.filter { $0.completed }.count
    }

    var todayRemainingCount: Int {
        max(habits.count - todayCompletedCount, 0)
    }

    func completedCount(_ window: Window) -> Int {
        let now = Date()
        let cal = Calendar.current
        return completions.filter { comp in
            switch window {
            case .day:
                return cal.isDate(comp.date, inSameDayAs: now)
            case .week:
                guard
                    let s1 = cal.dateInterval(of: .weekOfYear, for: comp.date),
                    let s2 = cal.dateInterval(of: .weekOfYear, for: now)
                else { return false }
                return s1 == s2
            case .month:
                return cal.component(.year, from: comp.date) == cal.component(.year, from: now)
                && cal.component(.month, from: comp.date) == cal.component(.month, from: now)
            }
        }.count
    }

    func toggleToday(for habit: Habit) {
        if let idx = completions.firstIndex(where: { $0.habitID == habit.id && Calendar.current.isDateInToday($0.date) }) {
            completions.remove(at: idx) // uncheck
        } else {
            completions.append(HabitCompletion(habitID: habit.id, date: Date()))
        }
        objectWillChange.send()
    }

    private func isCompletedToday(_ habit: Habit) -> Bool {
        completions.contains { $0.habitID == habit.id && Calendar.current.isDateInToday($0.date) }
    }

    func totalCapacity(_ window: Window) -> Int {
        let n = habits.count
        switch window {
        case .day:
            return n
        case .week:
            return n * daysInCurrentWeek
        case .month:
            return n * daysInCurrentMonth
        }
    }

    private var daysInCurrentWeek: Int { 7 }

    private var daysInCurrentMonth: Int {
        Calendar.current.range(of: .day, in: .month, for: Date())?.count ?? 30
    }

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

struct Habit: Identifiable, Hashable {
    let id: UUID = UUID()
    var name: String
}

struct HabitCompletion: Identifiable, Hashable {
    let id: UUID = UUID()
    let habitID: UUID
    let date: Date
}

struct HabitState: Identifiable {
    let id = UUID()
    let habit: Habit
    let completed: Bool
}

enum SampleData {
    static var seededCompletions: [HabitCompletion] {
        let cal = Calendar.current
        let now = Date()
        let daysAgo = { (d: Int) in cal.date(byAdding: .day, value: -d, to: now)! }

        // A few example completions in the past days/weeks for counters
        return [
            HabitCompletion(habitID: UUID(), date: daysAgo(20)), // random habit in month
            HabitCompletion(habitID: UUID(), date: daysAgo(6)),  // last week
            HabitCompletion(habitID: UUID(), date: daysAgo(1)),  // yesterday
        ]
    }
}

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

