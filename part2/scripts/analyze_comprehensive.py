import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import os
import glob
import re

def load_and_parse_results(base_dir):
    """Load both local and remote test results"""
    local_results = []
    remote_results = []
    
    # Parse local results
    for file in glob.glob(f"{base_dir}/optimization_tests_*/**.csv"):
        df = pd.read_csv(file)
        df['test_type'] = 'local'
        local_results.append(df)
    
    # Parse remote results
    for file in glob.glob(f"{base_dir}/remote_tests_*/**.csv"):
        df = pd.read_csv(file)
        df['test_type'] = 'remote'
        remote_results.append(df)
    
    return pd.concat(local_results + remote_results, ignore_index=True)

def analyze_bandwidth_limitations(df):
    """Analyze bandwidth limitations and bottlenecks"""
    # Calculate theoretical max bandwidth
    max_packet_size = df['packet_size'].max()
    theoretical_max = (max_packet_size * 1000) / (1024 * 1024)  # MB/s
    
    # Actual bandwidth statistics
    bandwidth_stats = df.groupby('test_type')['bandwidth'].agg(['mean', 'max', 'min'])
    
    # Calculate utilization
    bandwidth_stats['utilization'] = (bandwidth_stats['mean'] / theoretical_max) * 100
    
    return bandwidth_stats

def generate_comprehensive_report(results_dir):
    """Generate comprehensive performance report"""
    df = load_and_parse_results(results_dir)
    
    # Create report directory
    report_dir = os.path.join(results_dir, 'comprehensive_report')
    os.makedirs(report_dir, exist_ok=True)
    
    # Generate visualizations
    fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(15, 12))
    
    # 1. RTT Comparison (Local vs Remote)
    sns.boxplot(data=df, x='test_type', y='avg_rtt', hue='drop_rate', ax=ax1)
    ax1.set_title('RTT Distribution: Local vs Remote')
    ax1.set_ylabel('Average RTT (ms)')
    
    # 2. Bandwidth Comparison
    sns.barplot(data=df, x='packet_size', y='bandwidth', 
                hue='test_type', ax=ax2)
    ax2.set_title('Bandwidth by Packet Size')
    ax2.set_ylabel('Bandwidth (MB/s)')
    
    # 3. Packet Loss Analysis
    sns.scatterplot(data=df, x='drop_rate', y='final_loss',
                   hue='test_type', size='packet_size', ax=ax3)
    ax3.set_title('Packet Loss Analysis')
    ax3.set_ylabel('Final Loss Rate (%)')
    
    # 4. Optimization Impact
    sns.barplot(data=df[df['test_type']=='local'], 
                x='optimization', y='bandwidth',
                hue='packet_size', ax=ax4)
    ax4.set_title('Optimization Impact on Bandwidth')
    
    plt.tight_layout()
    plt.savefig(os.path.join(report_dir, 'comprehensive_analysis.png'))
    
    # Generate markdown report
    bandwidth_stats = analyze_bandwidth_limitations(df)
    
    with open(os.path.join(report_dir, 'performance_report.md'), 'w') as f:
        f.write("# UDP Performance Analysis Report\n\n")
        
        f.write("## Message Overhead\n")
        f.write(f"- Base message size: {df['packet_size'].min()} bytes\n")
        f.write(f"- Header overhead: 8 bytes (sequence number + checksum)\n")
        f.write(f"- ACK size: 4 bytes\n\n")
        
        f.write("## Round Trip Time (RTT)\n")
        f.write(f"- Local minimum RTT: {df[df['test_type']=='local']['avg_rtt'].min():.3f} ms\n")
        f.write(f"- Remote minimum RTT: {df[df['test_type']=='remote']['avg_rtt'].min():.3f} ms\n\n")
        
        f.write("## Bandwidth Analysis\n")
        f.write("### Local Testing\n")
        f.write(f"- Maximum: {bandwidth_stats.loc['local', 'max']:.2f} MB/s\n")
        f.write(f"- Average: {bandwidth_stats.loc['local', 'mean']:.2f} MB/s\n")
        f.write(f"- Utilization: {bandwidth_stats.loc['local', 'utilization']:.1f}%\n\n")
        
        f.write("### Remote Testing\n")
        f.write(f"- Maximum: {bandwidth_stats.loc['remote', 'max']:.2f} MB/s\n")
        f.write(f"- Average: {bandwidth_stats.loc['remote', 'mean']:.2f} MB/s\n")
        f.write(f"- Utilization: {bandwidth_stats.loc['remote', 'utilization']:.1f}%\n\n")
        
        f.write("## Performance Bottlenecks\n")
        f.write("1. Network latency (especially in remote testing)\n")
        f.write("2. Retry mechanism overhead\n")
        f.write("3. ACK waiting time\n")
        f.write("4. System call overhead\n\n")
        
        f.write("## Optimization Impact\n")
        opt_impact = df[df['test_type']=='local'].groupby('optimization')['bandwidth'].mean()
        f.write(f"- Optimized bandwidth: {opt_impact['optimized']:.2f} MB/s\n")
        f.write(f"- Unoptimized bandwidth: {opt_impact['unoptimized']:.2f} MB/s\n")
        f.write(f"- Performance improvement: {((opt_impact['optimized']/opt_impact['unoptimized'])-1)*100:.1f}%\n")

if __name__ == '__main__':
    results_dir = 'results'
    generate_comprehensive_report(results_dir)