#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#define MAX_SIZE 100
int Input(float (*)[2],int *);
int Calculate(float (*)[2], int,float *,float *,float *);
int Output(float,float,float,int);

int main(){
    float GPA[100][2];
    int success;
    int num_class;
    float total_Credit;

    success =Input(GPA,&num_class);//input the GPA into the matrix "GPA".
    if(success!=0){
        printf("Something wrong with function 'Input'!\n");
        return -1;
    }

    float averScore,averGPA;
    success=Calculate(GPA,num_class,&averScore,&averGPA,&total_Credit);
    if(success!=0){
        printf("Something wrong with function 'Calculate'!\n");
        return -1;
    }

    success =Output(averScore,averGPA,total_Credit,num_class);
    if(success!=0){
        printf("Something wrong with function 'Output'!\n");
        return -1;
    }
    printf("Enter any key to close this program...");
    getchar();
    return 0;
}

int Input(float (*GPA_point)[2], int *num_class){
    memset(GPA_point,0,sizeof(float)*200);
    printf("Waiting for reading files from GPA.txt...\n");
    FILE *fp;
    fp=fopen("GPA.csv","r");
    int i=0;
    char buffer[MAX_SIZE];
    fgets(buffer,MAX_SIZE,fp); // 读入的时候跳过第一行。
    while(fgets(buffer,MAX_SIZE,fp)!=NULL){
        char *result=NULL;
        result=strtok(buffer,",");
        int j=0;
        while(result!=NULL){
            if(j==4){//学分
               GPA_point[i][0]=atof(result);
            }
            else if(j==6){//分数
                GPA_point[i][1]=atof(result);
            }
            result=strtok(NULL,",");
            j++;
        }
        printf("credit is %f, score is:%f\n",GPA_point[i][0],GPA_point[i][1]);
        i++;
    }
    fclose(fp);
    printf("Done\n");
    *num_class=i;
    return 0;
}

int Calculate(float (*GPA_point)[2],int num_class,float *averScore,float *averGPA,float *total_Credit){
    *total_Credit=0;
    for(int i=0;i<num_class;i++){
        *total_Credit += GPA_point[i][0];
    }

    for(int i=0;i<num_class;i++){
        *averScore+=GPA_point[i][0]/(*total_Credit)*GPA_point[i][1];
        float GPA;
        if(GPA_point[i][1]>=85){
            GPA=4.0f;
        }else{
            GPA=4.0-(85-GPA_point[i][1])/10;
        }
        *averGPA+=GPA_point[i][0]/(*total_Credit)*GPA;
    }
    return 0;
}

int Output(float averScore, float averGPA,float total_Credit,int num_class){
    printf("Your number of classes is: %d\n", num_class);
    printf("Your total credit is: %f\n",total_Credit);
    printf("Your average Score is: %f\n",averScore);
    printf("Your average GPA is: %f\n",averGPA);
    return 0;
}