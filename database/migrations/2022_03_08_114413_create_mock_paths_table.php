<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('mock_paths', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('project_id')->index()->comment('项目id');
            $table->bigInteger('doc_id')->comment('文档id');
            $table->string('path')->comment('api路径');
            $table->string('format')->comment('数据格式');
            $table->string('method')->nullable()->comment('请求方法');
            $table->timestamps();

            $table->softDeletes();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('mock_paths');
    }
};
