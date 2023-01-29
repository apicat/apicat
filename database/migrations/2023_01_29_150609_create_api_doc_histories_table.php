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
        Schema::create('api_doc_histories', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('doc_id')->index()->comment('文档id');
            $table->string('title')->comment('文档名称');
            $table->mediumText('content')->nullable()->comment('文档内容');
            $table->bigInteger('last_user_id')->comment('最后编辑用户id');
            $table->timestamp('last_updated_at')->nullable()->comment('最后编辑时间');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('api_doc_histories');
    }
};
