import SwiftUI

struct HomeView: View {
    @StateObject private var vm = HabitViewModel()
    @State private var showSettings = false
    @State private var showTodaySheet = false
    @State private var dragOffset: CGFloat = 0

    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {
                VStack(alignment: .leading, spacing: 12) {
                    CounterLine(value: vm.statsMonth?.completed ?? 0, total: vm.statsMonth?.total ?? 0, suffix: "m", size: 84, opacity: 1.0)
                    CounterLine(value: vm.statsWeek?.completed  ?? 0, total: vm.statsWeek?.total  ?? 0, suffix: "w", size: 72, opacity: 0.72)
                    CounterLine(value: vm.statsDay?.completed   ?? 0, total: vm.statsDay?.total   ?? 0, suffix: "d", size: 60, opacity: 0.5)
                }
                .padding(.horizontal)
                .padding(.top, 12)
                .padding(.leading, 24)

                Spacer(minLength: 0)

                VStack(spacing: 12) {
                    Capsule().fill(Color.secondary.opacity(0.4)).frame(width: 40, height: 5).padding(.top, 8)
                    Button { showTodaySheet = true } label: {
                        HStack(spacing: 8) {
                            Image(systemName: "checkmark.circle.fill")
                            Text("\(vm.todayRemainingCount) remaining").font(.body)
                        }
                        .foregroundStyle(.secondary)
                        .padding(.vertical, 12)
                        .frame(maxWidth: .infinity)
                        .clipShape(RoundedRectangle(cornerRadius: 12))
                    }
                }
                .padding(.horizontal)
                .padding(.bottom, 24)
                .background(Color(.systemBackground).ignoresSafeArea(edges: .bottom).shadow(color: .black.opacity(0.05), radius: 8, y: -2))
                .gesture(DragGesture(minimumDistance: 10).onChanged { dragOffset = $0.translation.height }.onEnded { if $0.translation.height < -40 { showTodaySheet = true }; dragOffset = 0 })
            }
            .toolbar {
                ToolbarItem(placement: .principal) {
                    Menu {
                        ForEach(vm.competitions) { comp in
                            Button { vm.selectedCompetition = comp } label: { comp == vm.selectedCompetition ? AnyView(Label(comp.name, systemImage: "checkmark")) : AnyView(Text(comp.name)) }
                        }
                    } label: {
                        HStack(spacing: 4) { Text(vm.selectedCompetition.name).font(.body); Image(systemName: "chevron.down").font(.subheadline) }
                            .foregroundStyle(.secondary)
                            .padding(.vertical, 6)
                            .padding(.horizontal, 10)
                            .clipShape(Capsule())
                    }
                }
                ToolbarItem(placement: .topBarTrailing) { Button { showSettings = true } label: { Image(systemName: "gearshape").foregroundStyle(.secondary) } }
            }
        }
        .task { await vm.refreshAll() }
        .sheet(isPresented: $showSettings) { SettingsView().presentationDetents([.medium]) }
        .sheet(isPresented: $showTodaySheet) { TodayGoalsSheet(vm: vm).presentationDetents([.medium, .large]).presentationDragIndicator(.visible) }
    }
}

private struct CounterLine: View {
    let value: Int; let total: Int; let suffix: String; let size: CGFloat; let opacity: Double
    var body: some View { Text("\(value)/\(total)\(suffix)").font(.system(size: size, weight: .bold, design: .rounded)).foregroundStyle(Color.primary.opacity(opacity)).frame(maxWidth: .infinity, alignment: .leading) }
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
                            Button { Task { await vm.toggleToday(for: state.goal) } } label: {
                                Image(systemName: state.completed ? "checkmark.circle.fill" : "circle").foregroundStyle(state.completed ? .green : .secondary).imageScale(.large)
                            }.buttonStyle(.plain)
                        }
                    }
                }
            }
            .listStyle(.insetGrouped)
            .toolbar { ToolbarItem(placement: .principal) { Text("\(vm.todayCompletedCount)/\(vm.habits.count) completed").font(.subheadline).foregroundStyle(.secondary) } }
        }
    }
}


