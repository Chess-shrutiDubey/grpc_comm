import re
from typing import Dict, List, Tuple

def parse_bandwidth_results(content: str) -> Dict[int, float]:
    """Parse bandwidth test results from content."""
    results = {}
    pattern = r"Size: (\d+) bytes, Bandwidth: ([\d.]+) MB/s"
    matches = re.finditer(pattern, content)
    for match in matches:
        size = int(match.group(1))
        bandwidth = float(match.group(2))
        results[size] = bandwidth
    return results

def parse_marshal_results(content: str) -> List[Tuple[str, float, float]]:
    """Parse marshal test results from content."""
    results = []
    pattern = r"Result: {DataType:(\w+).*MarshalTime:([\d.]+)µs UnmarshalTime:([\d.]+)(?:µs|ns)"
    matches = re.finditer(pattern, content)
    for match in matches:
        data_type = match.group(1)
        marshal_time = float(match.group(2))
        unmarshal_time = float(match.group(3))
        if 'ns' in match.group(0):
            unmarshal_time /= 1000  # Convert ns to µs
        results.append((data_type, marshal_time, unmarshal_time))
    return results

def parse_rtt_results(content: str) -> Dict[str, float]:
    """Parse RTT test results from content."""
    results = {}
    patterns = {
        'first': r"First RTT: ([\d.]+)(?:ms|µs)",
        'min': r"Min RTT: ([\d.]+)(?:ms|µs)",
        'max': r"Max RTT: ([\d.]+)(?:ms|µs)",
        'avg': r"Avg RTT: ([\d.]+)(?:ms|µs)"
    }
    
    for key, pattern in patterns.items():
        match = re.search(pattern, content)
        if match:
            value = float(match.group(1))
            if 'ms' in match.group(0):
                value *= 1000  # Convert ms to µs
            results[key] = value
    return results

def analyze_optimization():
    # Read results files
    with open('results/unoptimized_results.txt', 'r') as f:
        unopt = f.read()
    with open('results/optimized_results.txt', 'r') as f:
        opt = f.read()

    # Parse results
    unopt_bandwidth = parse_bandwidth_results(unopt)
    opt_bandwidth = parse_bandwidth_results(opt)
    unopt_marshal = parse_marshal_results(unopt)
    opt_marshal = parse_marshal_results(opt)
    unopt_rtt = parse_rtt_results(unopt)
    opt_rtt = parse_rtt_results(opt)

    # Generate comparison CSV
    with open('results/optimization_analysis.csv', 'w') as f:
        # Write header
        f.write("Test Type,Metric,Unoptimized,Optimized,Difference (%)\n")
        
        # Bandwidth comparisons
        for size in sorted(unopt_bandwidth.keys()):
            unopt_val = unopt_bandwidth[size]
            opt_val = opt_bandwidth[size]
            diff_percent = ((opt_val - unopt_val) / unopt_val) * 100
            f.write(f"Bandwidth {size}B,MB/s,{unopt_val:.2f},{opt_val:.2f},{diff_percent:.2f}\n")

        # Marshal comparisons
        for (unopt_type, unopt_m, unopt_um), (opt_type, opt_m, opt_um) in zip(unopt_marshal, opt_marshal):
            f.write(f"Marshal {unopt_type},Marshal Time (µs),{unopt_m:.3f},{opt_m:.3f},{((opt_m - unopt_m) / unopt_m) * 100:.2f}\n")
            f.write(f"Marshal {unopt_type},Unmarshal Time (µs),{unopt_um:.3f},{opt_um:.3f},{((opt_um - unopt_um) / unopt_um) * 100:.2f}\n")

        # RTT comparisons
        for metric in ['first', 'min', 'max', 'avg']:
            unopt_val = unopt_rtt[metric]
            opt_val = opt_rtt[metric]
            diff_percent = ((opt_val - unopt_val) / unopt_val) * 100
            f.write(f"RTT,{metric.capitalize()} (µs),{unopt_val:.3f},{opt_val:.3f},{diff_percent:.2f}\n")

if __name__ == "__main__":
    analyze_optimization()
    print("Optimization analysis complete. Results written to results/optimization_analysis.csv")