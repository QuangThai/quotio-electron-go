import {
  Bar,
  BarChart,
  Cell,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { cn } from "../../lib/utils";

// Neobrutalism color palette
const COLORS = [
  "#2563EB",
  "#059669",
  "#7C3AED",
  "#EA580C",
  "#DB2777",
  "#D97706",
  "#0891B2",
  "#4F46E5",
];

interface ProviderChartData {
  name: string;
  value: number;
  accounts: number;
  icon?: string;
  [key: string]: string | number | undefined;
}

interface QuotaChartData {
  name: string;
  used: number;
  limit: number;
  percentage: number;
}

interface ProviderDistributionChartProps {
  data: ProviderChartData[];
  title?: string;
  className?: string;
}

export function ProviderDistributionChart({
  data,
  title = "Provider Distribution",
  className,
}: ProviderDistributionChartProps) {
  if (data.length === 0) {
    return (
      <div
        className={cn(
          "bg-gray-50 border-2 border-dashed border-gray-300 rounded-sm p-8 text-center",
          className
        )}
      >
        <p className="text-gray-500 text-sm">No providers to display</p>
      </div>
    );
  }

  return (
    <div
      className={cn(
        "bg-white border-2 border-black rounded-sm shadow-neobrutal p-4",
        className
      )}
    >
      <h4 className="font-bold text-sm mb-4 flex items-center gap-2">
        <span className="w-3 h-3 bg-neobrutal-blue rounded-full"></span>
        {title}
      </h4>
      <div className="h-48">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="50%"
              innerRadius={40}
              outerRadius={70}
              paddingAngle={2}
              dataKey="accounts"
              stroke="#000"
              strokeWidth={2}
            >
              {data.map((_, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={COLORS[index % COLORS.length]}
                />
              ))}
            </Pie>
            <Tooltip
              contentStyle={{
                backgroundColor: "#fff",
                border: "2px solid #000",
                borderRadius: "4px",
                boxShadow: "3px 3px 0 #000",
              }}
              formatter={(value) => [`${value} accounts`, "Accounts"]}
            />
          </PieChart>
        </ResponsiveContainer>
      </div>
      <div className="flex flex-wrap gap-2 mt-4 justify-center">
        {data.map((item, index) => (
          <div key={item.name} className="flex items-center gap-1.5 text-xs">
            <span
              className="w-3 h-3 rounded-sm border border-black"
              style={{ backgroundColor: COLORS[index % COLORS.length] }}
            />
            <span className="font-medium">{item.name}</span>
            <span className="text-gray-500">({item.accounts})</span>
          </div>
        ))}
      </div>
    </div>
  );
}

interface QuotaUsageChartProps {
  data: QuotaChartData[];
  title?: string;
  className?: string;
}

export function QuotaUsageChart({
  data,
  title = "Quota Usage",
  className,
}: QuotaUsageChartProps) {
  if (data.length === 0) {
    return (
      <div
        className={cn(
          "bg-gray-50 border-2 border-dashed border-gray-300 rounded-sm p-8 text-center",
          className
        )}
      >
        <p className="text-gray-500 text-sm">No quota data to display</p>
      </div>
    );
  }

  return (
    <div
      className={cn(
        "bg-white border-2 border-black rounded-sm shadow-neobrutal p-4",
        className
      )}
    >
      <h4 className="font-bold text-sm mb-4 flex items-center gap-2">
        <span className="w-3 h-3 bg-neobrutal-green rounded-full"></span>
        {title}
      </h4>
      <div className="h-48">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={data}
            layout="vertical"
            margin={{ left: 0, right: 20 }}
          >
            <XAxis type="number" hide />
            <YAxis
              dataKey="name"
              type="category"
              width={80}
              tick={{ fontSize: 11, fontWeight: 600 }}
              tickLine={false}
              axisLine={false}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: "#fff",
                border: "2px solid #000",
                borderRadius: "4px",
                boxShadow: "3px 3px 0 #000",
              }}
              formatter={(value) => {
                const numValue = Number(value);
                return [`${numValue.toFixed(1)}%`, "Usage"];
              }}
            />
            <Bar
              dataKey="percentage"
              fill="#2563EB"
              radius={[0, 4, 4, 0]}
              stroke="#000"
              strokeWidth={1}
              maxBarSize={24}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

interface StatsCardProps {
  label: string;
  value: string | number;
  icon?: React.ReactNode;
  color?: "blue" | "green" | "purple" | "orange";
  trend?: "up" | "down" | "neutral";
  className?: string;
}

export function StatsCard({
  label,
  value,
  icon,
  color = "blue",
  className,
}: StatsCardProps) {
  const colorClasses = {
    blue: "bg-neobrutal-blue/10 text-neobrutal-blue border-neobrutal-blue/20",
    green:
      "bg-neobrutal-green/10 text-neobrutal-green border-neobrutal-green/20",
    purple:
      "bg-neobrutal-purple/10 text-neobrutal-purple border-neobrutal-purple/20",
    orange:
      "bg-neobrutal-orange/10 text-neobrutal-orange border-neobrutal-orange/20",
  };

  return (
    <div
      className={cn(
        "bg-white border-2 border-black rounded-base shadow-neobrutal-sm p-4",
        "hover:shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px]",
        "transition-all duration-200",
        className
      )}
    >
      <div className="flex items-center justify-between mb-2">
        <span className="text-[10px] font-bold text-gray-500 uppercase tracking-wider">
          {label}
        </span>
        {icon && (
          <div
            className={cn(
              "p-1.5 rounded-base border-2 border-black/10",
              colorClasses[color]
            )}
          >
            {icon}
          </div>
        )}
      </div>
      <div className="text-xl font-black">{value}</div>
    </div>
  );
}
