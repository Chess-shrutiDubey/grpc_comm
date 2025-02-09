import json
import csv
import os
import glob
from datetime import datetime
import statistics

def parse_duration(duration_str):
    """Convert duration string to milliseconds"""
    if 'ms' in duration_str:
        return float(duration_str.replace('ms', ''))
    elif 'µs' in duration_str:
        return float(duration_str.replace('µs', '')) / 1000
    return 0

def process_performance_results(root_dir):
    all_results = {
        'rtt': [],
        'bandwidth': [],
        'marshal': []
    }
    
    # Processings each test type with detailed debugging
    for test_type in ['rtt', 'bandwidth', 'marshal']:
        json_dir = os.path.join(root_dir, 'results', test_type)
        if os.path.exists(json_dir):
            json_files = glob.glob(os.path.join(json_dir, '*.json'))
            print(f"\nFound {len(json_files)} {test_type} files in {json_dir}")
            
            for json_file in json_files:
                try:
                    with open(json_file, 'r') as f:
                        data = json.load(f)
                        all_results[test_type].append(data)
                except json.JSONDecodeError as e:
                    print(f"Error reading {json_file}: {e}")
                except Exception as e:
                    print(f"Unexpected error with {json_file}: {e}")

    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    
    # RTT Results
    if all_results['rtt']:
        csv_file = os.path.join(root_dir, 'results', f'rtt_results_{timestamp}.csv')
        with open(csv_file, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(['Size (bytes)', 'Average RTT (ms)', 'Min RTT (ms)', 'Max RTT (ms)'])
            for result in all_results['rtt']:
                try:
                    if 'rtts' in result and result['rtts']:  # Check if rtts exists and is not empty
                        # Convert all RTT values to milliseconds
                        rtts_ms = [parse_duration(rtt) for rtt in result['rtts']]
                        if rtts_ms:  # Check if we have valid measurements
                            writer.writerow([
                                result['message_size'],
                                f"{statistics.mean(rtts_ms):.3f}",
                                f"{min(rtts_ms):.3f}",
                                f"{max(rtts_ms):.3f}"
                            ])
                            print(f"Processed RTT data for size {result['message_size']}")
                        else:
                            print(f"Warning: No valid RTT measurements for size {result['message_size']}")
                    else:
                        print(f"Warning: Missing RTT data in result: {result.keys()}")
                except Exception as e:
                    print(f"Error processing RTT result: {e}")
                    print(f"Problematic result: {result}")

    # Bandwidth Results
    if all_results['bandwidth']:
        csv_file = os.path.join(root_dir, 'results', f'bandwidth_results_{timestamp}.csv')
        with open(csv_file, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(['Size (bytes)', 'Bandwidth (MB/s)', 'Timestamp'])
            for result in all_results['bandwidth']:
                try:
                    if all(key in result for key in ['message_size', 'bandwidth_mbps']):
                        writer.writerow([
                            result['message_size'],
                            result['bandwidth_mbps'],
                            result.get('timestamp', '')
                        ])
                        print(f"Processed bandwidth data for size {result['message_size']}")
                    else:
                        print(f"Warning: Missing required bandwidth data in result: {result.keys()}")
                except Exception as e:
                    print(f"Error processing bandwidth result: {e}")
                    print(f"Problematic result: {result}")

    # Marshal Results
    if all_results['marshal']:
        csv_file = os.path.join(root_dir, 'results', f'marshal_results_{timestamp}.csv')
        with open(csv_file, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(['Size (bytes)', 'Marshal Time (ms)', 'Data Type', 'Timestamp'])
            for result in all_results['marshal']:
                writer.writerow([
                    result.get('message_size'),
                    result.get('marshal_time'),
                    result.get('data_type'),
                    result.get('timestamp')
                ])

def main():
    project_root = "/home/shruti/SEM2/grpc_comm/grpc_comm/part3"
    process_performance_results(project_root)
    print("Results have been converted to CSV format")

if __name__ == "__main__":
    main()