<?php

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldBeUnique;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use App\Repositories\Export\PostmanRepository;

class PostmanExport implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /**
     * 导出任务ID
     *
     * @var int
     */
    public $jobID;

    /**
     * Create a new job instance.
     *
     * @param string $jobID 导出任务ID
     * @return void
     */
    public function __construct($jobID)
    {
        $this->jobID = $jobID;
    }

    /**
     * Execute the job.
     *
     * @return void
     */
    public function handle()
    {
        (new PostmanRepository)->startJob($this->jobID);
    }
}
