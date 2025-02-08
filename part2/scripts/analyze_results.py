import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import glob
import os
import re

def parse_csv_file(file_path):
    try:
        # Read CSV file, skipping any non-CSV output
        with open(file_path, 'r') as f:
            lines = f.readlines()
        
        # Find the start of CSV data
        start_idx = -1
        for i, line in enumerate(lines):
            if line.strip() == "Metric,Value":
                start_idx = i
                break
        
        if start_idx == -1:
            print(f"Warning: No CSV header found in {file_path}")
            return None
        
        # Parse CSV data
        metrics = {}
        for line in lines[start_idx+1:]:
            if ',' not in line:
                continue
            metric, value = line.strip().split(',')
            try:
                metrics[metric] = float(value)
            except ValueError:
                print(f"Warning: Could not parse value for {metric} in {file_path}")
                continue
        
        # Map metrics to results
        results = {
            'packets_sent': int(metrics.get('Packets_Sent', 0)),
            'packets_received': int(metrics.get('Packets_Received', 0)),
            'dropped_packets': int(metrics.get('Dropped_Packets', 0)),
            'packet_loss': float(metrics.get('Packet_Loss_Rate', 0)),
            'bandwidth': float(metrics.get('Bandwidth_MBps', 0)),
            'avg_rtt': float(metrics.get('Average_RTT_ms', 0))
        }
        
        return pd.Series(results)
    except Exception as e:
        print(f"Error parsing {file_path}: {str(e)}")
        return None

def load_results(results_dir):
    data = []
    
    # Print all available files for debugging
    print("\nAvailable files in directory:")
    for file in glob.glob(os.path.join(results_dir, "*.csv")):
        print(os.path.basename(file))
    
    # Process both optimized and unoptimized files
    for opt_type in ["optimized", "unoptimized"]:
        file_pattern = os.path.join(results_dir, f"{opt_type}_rate*_size*.csv")
        for file in glob.glob(file_pattern):
            try:
                filename = os.path.basename(file)
                rate = int(re.search(r'rate(\d+)', filename).group(1))
                size = int(re.search(r'size(\d+)', filename).group(1))
                
                results = parse_csv_file(file)
                if results is not None:
                    results['optimization'] = opt_type
                    results['drop_rate'] = rate
                    results['packet_size'] = size
                    data.append(results)
                    print(f"Processed {opt_type} file: {filename}")
            except Exception as e:
                print(f"Error processing {file}: {e}")
    
    if not data:
        raise ValueError("No valid data files found")
    
    # Create DataFrame and print summary
    df = pd.DataFrame(data)
    print("\nData loading summary:")
    print(f"Total files processed: {len(data)}")
    print("Files per optimization type:")
    print(df['optimization'].value_counts())
    
    return df

def analyze_performance():
    results_dirs = glob.glob("../results/optimization_tests_*")
    if not results_dirs:
        print("No results found. Run optimization tests first")
        return

    try:
        latest_dir = max(results_dirs, key=os.path.getctime)
        print(f"Analyzing results from: {latest_dir}")
        df = load_results(latest_dir)
        
        # Set custom color palette for optimization types
        custom_palette = {'optimized': '#2ecc71', 'unoptimized': '#e67e22'}
        sns.set_style("whitegrid")
        
        # Create plots with larger figure size
        fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(16, 12))
        
        # 1. RTT vs Drop Rate
        sns.boxplot(data=df, x='drop_rate', y='avg_rtt', hue='optimization', 
                   palette=custom_palette, ax=ax1)
        ax1.set_title('Average RTT vs Drop Rate', fontsize=12, pad=10)
        ax1.set_ylabel('Average RTT (ms)', fontsize=10)
        ax1.set_xlabel('Drop Rate (%)', fontsize=10)
        ax1.legend(title='Optimization', title_fontsize=10, fontsize=9)

        # 2. Bandwidth vs Packet Size
        sns.barplot(data=df, x='packet_size', y='bandwidth', hue='optimization', 
                   palette=custom_palette, ax=ax2)
        ax2.set_title('Bandwidth vs Packet Size', fontsize=12, pad=10)
        ax2.set_ylabel('Bandwidth (MB/s)', fontsize=10)
        ax2.set_xlabel('Packet Size (bytes)', fontsize=10)
        ax2.legend(title='Optimization', title_fontsize=10, fontsize=9)

        # 3. Dropped Packets vs Drop Rate
        sns.barplot(data=df, x='drop_rate', y='dropped_packets', hue='optimization',
                    errorbar=('ci', 95), palette=custom_palette, ax=ax3)
        ax3.set_title('Dropped Packets vs Drop Rate', fontsize=12, pad=10)
        ax3.set_ylabel('Number of Dropped Packets', fontsize=10)
        ax3.set_xlabel('Drop Rate (%)', fontsize=10)
        ax3.legend(title='Optimization', title_fontsize=10, fontsize=9)

        # 4. Packet Loss vs Drop Rate
        sns.lineplot(data=df, x='drop_rate', y='packet_loss', hue='optimization',
                    errorbar=('ci', 95), err_style='band', marker='o', 
                    palette=custom_palette, ax=ax4)
        ax4.set_title('Packet Loss Rate vs Drop Rate', fontsize=12, pad=10)
        ax4.set_ylabel('Packet Loss Rate (%)', fontsize=10)
        ax4.set_xlabel('Drop Rate (%)', fontsize=10)
        ax4.legend(title='Optimization', title_fontsize=10, fontsize=9)

        # Adjust layout to prevent overlap
        plt.tight_layout(pad=3.0)
        
        # Save plot with higher DPI for better quality
        output_file = os.path.join(latest_dir, 'performance_analysis.png')
        plt.savefig(output_file, dpi=300, bbox_inches='tight')
        print(f"Analysis saved to {output_file}")

        # Print full summary statistics
        pd.set_option('display.max_columns', None)
        pd.set_option('display.width', None)
        print("\nPerformance Summary:")
        summary = df.groupby(['optimization', 'drop_rate']).agg({
            'avg_rtt': 'mean',
            'bandwidth': 'mean',
            'dropped_packets': 'mean',
            'packet_loss': 'mean'
        }).round(3)
        print(summary.to_string())
        
        # Save detailed summary
        summary.to_csv(os.path.join(latest_dir, 'summary.csv'))

    except Exception as e:
        print(f"Analysis failed: {e}")

if __name__ == '__main__':
    analyze_performance()