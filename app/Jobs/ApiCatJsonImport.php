<?php

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldBeUnique;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use App\Repositories\Import\ApiCatRepository;

class ApiCatJsonImport implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /**
     * 任务可尝试的次数
     *
     * @var int
     */
    public $tries = 2;

    /**
     * 导入任务ID
     *
     * @var int
     */
    public $jobID;

    /**
     * Create a new job instance.
     *
     * @param string $jobID 导入任务ID
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
        (new ApiCatRepository)->startJob($this->jobID);
    }
}
