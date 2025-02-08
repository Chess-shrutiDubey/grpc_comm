import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import glob
import os
import re

def parse_log_file(file_path):
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Look for the Results section specifically
    results_section = content.split('Results:')[-1] if 'Results:' in content else content
    
    results = {}
    try:
        patterns = {
            'packets_sent': r'Packets sent: (\d+)',
            'packets_received': r'Packets received: (\d+)',
            'initial_drops': r'Initial drops: (\d+)',
            'final_loss': r'Final packet loss: ([\d.]+)%',
            'bandwidth': r'Bandwidth: ([\d.]+)',
            'avg_rtt': r'Average RTT: ([\d.]+)ms',
            'duration': r'Total duration: ([\d.]+)s'
        }
        
        for key, pattern in patterns.items():
            match = re.search(pattern, results_section)
            if not match:
                print(f"Warning: Could not find {key} in {file_path}")
                return None
            results[key] = float(match.group(1))
            
        return pd.Series(results)
    except Exception as e:
        print(f"Error parsing {file_path}: {e}")
        return None

def load_results(results_dir):
    seen_files = set()  # Track processed files
    data = []
    
    # Print all available files for debugging
    print("\nAvailable files in directory:")
    for file in glob.glob(os.path.join(results_dir, "*.csv")):
        print(os.path.basename(file))
    
    # Process both optimized and unoptimized files
    for opt_type in ["optimized", "unoptimized"]:
        file_pattern = os.path.join(results_dir, f"{opt_type}_rate*_size*.csv")
        for file in glob.glob(file_pattern):
            if file in seen_files:
                continue
            seen_files.add(file)
            
            try:
                filename = os.path.basename(file)
                rate = int(re.search(r'rate(\d+)', filename).group(1))
                size = int(re.search(r'size(\d+)', filename).group(1))
                
                results = parse_log_file(file)
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

        # 3. Initial Drops vs Drop Rate
        sns.barplot(data=df, x='drop_rate', y='initial_drops', hue='optimization', 
                   palette=custom_palette, ax=ax3)
        ax3.set_title('Initial Drops vs Drop Rate', fontsize=12, pad=10)
        ax3.set_ylabel('Initial Drops', fontsize=10)
        ax3.set_xlabel('Drop Rate (%)', fontsize=10)
        ax3.legend(title='Optimization', title_fontsize=10, fontsize=9)

        # 4. Final Loss Rate vs Drop Rate
        sns.lineplot(data=df, x='drop_rate', y='final_loss', hue='optimization', 
                    palette=custom_palette, marker='o', markersize=8, ax=ax4)
        ax4.set_title('Final Loss Rate vs Drop Rate', fontsize=12, pad=10)
        ax4.set_ylabel('Final Loss Rate (%)', fontsize=10)
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
            'initial_drops': 'mean',
            'final_loss': 'mean'
        }).round(3)
        print(summary.to_string())
        
        # Save detailed summary
        summary.to_csv(os.path.join(latest_dir, 'summary.csv'))

    except Exception as e:
        print(f"Analysis failed: {e}")

if __name__ == '__main__':
    analyze_performance()